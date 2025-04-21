package api

import (
	"context"
	"encoding/json"

	"net/http"
	"time"

	"backend-go/config"
	"backend-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func init() {
	client = config.ConnectDB()
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message

	// Decodificăm corpul cererii
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Setăm ID-ul și timestamp-ul
	message.ID = primitive.NewObjectID()
	message.Timestamp = time.Now()

	// Inserăm mesajul în MongoDB
	collection := config.GetCollection(client, "messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Răspundem cu mesajul creat
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
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
