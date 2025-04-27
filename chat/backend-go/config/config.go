package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	// Citim MONGO_URI din variabila de mediu
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("Eroare: Variabila de mediu MONGO_URI nu este setată")
	}

	// Configurăm conexiunea la MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("Eroare la crearea clientului MongoDB:", err)
	}

	// Setăm un timeout pentru conexiune
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Eroare la conectarea la MongoDB:", err)
	}

	// Verificăm conexiunea
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Eroare la ping MongoDB:", err)
	}

	log.Println("Conectat la MongoDB!")
	return client
}

// Funcție pentru a obține colecția de mesaje
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// Citim numele bazei de date din variabila de mediu sau folosim o valoare implicită
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "chatdb" // Valoare implicită
		log.Println("Variabila DB_NAME nu este setată, se folosește valoarea implicită:", dbName)
	}

	return client.Database(dbName).Collection(collectionName)
}
