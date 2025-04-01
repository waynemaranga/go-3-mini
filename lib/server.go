package lib

import (
	"encoding/json"
	"net/http"
)

type ChatRequest struct {
	Message string `json:"message"`
}

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

func StartServer(port string) {
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/prompt", promptHandler)
	http.HandleFunc("/chats", chatsHandler)

	LogInfo("Server starting on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		LogError(err)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get chat history
	history := GetChatHistory()

	// Add new user message
	userMessage := ChatMessage{Role: "user", Content: req.Message}
	history = append(history, userMessage)
	SaveChat(userMessage)

	// Get AI response
	aiResponse := GetAIResponse(history)
	aiMessage := ChatMessage{Role: "assistant", Content: aiResponse}
	SaveChat(aiMessage)

	// Return response
	response := map[string]string{"response": aiResponse}
	json.NewEncoder(w).Encode(response)
}

func promptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create single message conversation
	messages := []ChatMessage{
		{Role: "user", Content: req.Prompt},
	}

	// Get AI response
	response := GetAIResponse(messages)

	// Return response
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func chatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history := GetChatHistory()
	json.NewEncoder(w).Encode(history)
}
