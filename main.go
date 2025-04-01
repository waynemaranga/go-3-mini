// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go-3-mini/db"
	// "go-3-mini/clients"
	"mini-chat/openai"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present
	godotenv.Load()

	// Get configuration from environment
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := db.Connect(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB!")

	// Initialize OpenAI client
	openaiClient := openai.NewClient(openAIAPIKey)
	if openAIAPIKey == "" {
		log.Println("Warning: OPENAI_API_KEY not set. AI features will not work.")
	}

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(mongoClient, openaiClient)

	// Set up HTTP routes
	http.HandleFunc("/chat", chatHandler.HandleChat)
	http.HandleFunc("/chats", chatHandler.GetChats)
	http.HandleFunc("/ai-chat", chatHandler.HandleAIChat)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
