package main

import (
	"backend-go/config"
	"backend-go/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	clients   = make(map[*websocket.Conn]bool) // Lista de clienți conectați
	broadcast = make(chan models.Message)      // Canal pentru mesaje broadcast
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Permitem orice origine în scopuri de dezvoltare
		},
	}
	client *mongo.Client
)

func init() {
	client = config.ConnectDB()
}

func main() {
	// Citim portul din variabila de mediu
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Valoare implicită
		log.Println("Variabila SERVER_PORT nu este setată, se folosește portul implicit:", port)
	}

	router := mux.NewRouter()

	// Rută pentru WebSocket
	router.HandleFunc("/ws", handleConnections)

	router.HandleFunc("/api/messages", getMessages).Methods("GET")

	// Adaugăm middleware-ul CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Permitem cereri de la orice origine
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(router)

	// Pornim un goroutine pentru gestionarea mesajelor
	go handleMessages()

	log.Println("Server pornit pe port", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade-ăm conexiunea HTTP la WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Înregistrăm clientul nou
	clients[ws] = true

	// Trimitem mesajele anterioare din MongoDB
	sendPreviousMessages(ws)

	// Bucla infinită pentru a menține conexiunea deschisă
	for {
		var msg models.Message
		// Citim mesajul nou de la WebSocket
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Eroare: %v", err)
			delete(clients, ws)
			break
		}

		// Procesăm mesajul primit
		msg.ID = primitive.NewObjectID()
		msg.Timestamp = time.Now()

		// Salvăm mesajul în MongoDB
		collection := config.GetCollection(client, "messages")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = collection.InsertOne(ctx, msg)
		cancel()
		if err != nil {
			log.Printf("Eroare la salvarea în MongoDB: %v", err)
		}

		// Trimitem mesajul către toți clienții
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Primim mesaj de la canalul de broadcast
		msg := <-broadcast
		// Trimitem către toți clienții conectați
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Eroare: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func sendPreviousMessages(conn *websocket.Conn) {
	var messages []models.Message

	// Obținem mesajele din MongoDB
	collection := config.GetCollection(client, "messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Eroare la obținerea mesajelor: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Iterăm prin rezultate
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			log.Printf("Eroare la decodarea mesajului: %v", err)
			continue
		}
		messages = append(messages, message)
	}

	// Trimitem fiecare mesaj către client
	for _, msg := range messages {
		conn.WriteJSON(msg)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	var messages []models.Message

	// Obținem mesajele din MongoDB
	collection := config.GetCollection(client, "messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Iterăm prin rezultate
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	// Răspundem cu lista de mesaje
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
