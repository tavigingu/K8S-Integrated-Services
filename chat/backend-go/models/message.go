package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `json:"id"`
	Content   string             `json:"content"`
	Username  string             `json:"username"`
	Timestamp time.Time          `json:"timestamp"`
}
