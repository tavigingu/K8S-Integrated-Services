package dbmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id"`
	Content   string             `bson:"content"`
	Username  string             `bson:"username"`
	Timestamp time.Time          `bson:"timestamp"`
}
