package db

import (
	"context"
	"fmt"
	"go-3-mini/config"
	"go-3-mini/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func ConnectDB() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database(config.DBName).Collection(config.Collection)
	fmt.Println("Connected to MongoDB")
}

func SaveChat(chat models.ChatMessage) {
	_, err := collection.InsertOne(context.TODO(), chat)
	if err != nil {
		log.Println("Error saving chat:", err)
	}
}

func GetChatHistory() []models.ChatMessage {
	var chats []models.ChatMessage
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error fetching chats:", err)
		return chats
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var chat models.ChatMessage
		cursor.Decode(&chat)
		chats = append(chats, chat)
	}
	return chats
}
