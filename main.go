package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Message represents a single chat message
type Message struct {
	User      string    `json:"user" bson:"user"`
	Content   string    `json:"content" bson:"content"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// Chat represents a collection of messages
type Chat struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Messages  []Message          `json:"messages" bson:"messages"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// OpenAI API related structures
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string         `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var (
	client         *mongo.Client
	chatCollection *mongo.Collection
	openAIAPIKey   string
)

func main() {
	// Get OpenAI API key from environment
	openAIAPIKey = os.Getenv("OPENAI_API_KEY")
	if openAIAPIKey == "" {
		log.Println("Warning: OPENAI_API_KEY not set. AI features will not work.")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	
	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	
	log.Println("Connected to MongoDB!")
	chatCollection = client.Database("mini-chat").Collection("chats")
	
	// Set up HTTP routes
	http.HandleFunc("/chat", handleChat)
	http.HandleFunc("/chats", getChats)
	http.HandleFunc("/ai-chat", handleAIChat)
	
	// Start the server
	log.Println("Server starting on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Set timestamp if not provided
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	
	// Create a new chat or append to existing
	chatID := r.URL.Query().Get("chatId")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if chatID == "" {
		// Create new chat
		chat := Chat{
			Title:     "New Chat",
			Messages:  []Message{msg},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		result, err := chatCollection.InsertOne(ctx, chat)
		if err != nil {
			http.Error(w, "Failed to save chat", http.StatusInternalServerError)
			log.Printf("Error saving chat: %v", err)
			return
		}
		
		// Get the inserted ID and send back to client
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			chat.ID = oid
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chat)
	} else {
		// Try to parse the ID
		objID, err := primitive.ObjectIDFromHex(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		
		// Append message to existing chat
		filter := bson.M{"_id": objID}
		update := bson.M{
			"$push": bson.M{"messages": msg},
			"$set":  bson.M{"updatedAt": time.Now()},
		}
		
		result, err := chatCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(w, "Failed to update chat", http.StatusInternalServerError)
			log.Printf("Error updating chat: %v", err)
			return
		}
		
		if result.MatchedCount == 0 {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "Message added to chat"}`))
	}
}

func getChats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cursor, err := chatCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error retrieving chats", http.StatusInternalServerError)
		log.Printf("Error retrieving chats: %v", err)
		return
	}
	defer cursor.Close(ctx)
	
	var chats []Chat
	if err := cursor.All(ctx, &chats); err != nil {
		http.Error(w, "Error parsing chats", http.StatusInternalServerError)
		log.Printf("Error parsing chats: %v", err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func handleAIChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if openAIAPIKey == "" {
		http.Error(w, "OpenAI API key not configured", http.StatusInternalServerError)
		return
	}
	
	var userMsg Message
	err := json.NewDecoder(r.Body).Decode(&userMsg)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Set timestamp if not provided
	if userMsg.Timestamp.IsZero() {
		userMsg.Timestamp = time.Now()
	}
	
	// Get existing chat or create new one
	chatID := r.URL.Query().Get("chatId")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var chat Chat
	var messages []Message
	
	if chatID != "" {
		// Try to find existing chat
		objID, err := primitive.ObjectIDFromHex(chatID)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		
		err = chatCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&chat)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Chat not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
				log.Printf("Error finding chat: %v", err)
			}
			return
		}
		
		messages = chat.Messages
	} else {
		// New chat
		chat = Chat{
			Title:     "AI Chat",
			Messages:  []Message{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	
	// Add user message to the chat
	messages = append(messages, userMsg)
	
	// Prepare messages for OpenAI API
	openAIMessages := []OpenAIMessage{}
	for _, msg := range messages {
		role := "user"
		if msg.User == "ai" {
			role = "assistant"
		}
		openAIMessages = append(openAIMessages, OpenAIMessage{
			Role:    role,
			Content: msg.Content,
		})
	}
	
	// Call OpenAI API
	aiResponse, err := callOpenAI(openAIMessages)
	if err != nil {
		http.Error(w, "Error calling OpenAI API", http.StatusInternalServerError)
		log.Printf("OpenAI API error: %v", err)
		return
	}
	
	// Create AI message
	aiMsg := Message{
		User:      "ai",
		Content:   aiResponse,
		Timestamp: time.Now(),
	}
	
	// Add AI response to messages
	messages = append(messages, aiMsg)
	
	// Save or update the chat
	if chatID == "" {
		// It's a new chat
		chat.Messages = messages
		result, err := chatCollection.InsertOne(ctx, chat)
		if err != nil {
			http.Error(w, "Failed to save chat", http.StatusInternalServerError)
			log.Printf("Error saving chat: %v", err)
			return
		}
		
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			chat.ID = oid
		}
	} else {
		// Update existing chat
		objID, _ := primitive.ObjectIDFromHex(chatID)
		_, err = chatCollection.UpdateOne(
			ctx,
			bson.M{"_id": objID},
			bson.M{
				"$set": bson.M{
					"messages":  messages,
					"updatedAt": time.Now(),
				},
			},
		)
		if err != nil {
			http.Error(w, "Failed to update chat", http.StatusInternalServerError)
			log.Printf("Error updating chat: %v", err)
			return
		}
		
		chat.Messages = messages
	}
	
	// Return the updated chat
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

func callOpenAI(messages []OpenAIMessage) (string, error) {
	reqBody := OpenAIRequest{
		Model:    "o3-mini", // Specify o3 mini model
		Messages: messages,
	}
	
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}
	
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIAPIKey)
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	
	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}
	
	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}
	
	return openAIResp.Choices[0].Message.Content, nil
}