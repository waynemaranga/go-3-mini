package lib

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func ConnectDB() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database(DBName).Collection(Collection)
	fmt.Println("✅ Connected to MongoDB")
}

func SaveChat(chat ChatMessage) {
	_, err := collection.InsertOne(context.TODO(), chat)
	if err != nil {
		log.Println("⛔ Error saving chat:", err)
	}
}

func GetChatHistory() []ChatMessage {
	var chats []ChatMessage
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("⛔ Error fetching chats:", err)
		return chats
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var chat ChatMessage
		cursor.Decode(&chat)
		chats = append(chats, chat)
	}
	return chats
}
