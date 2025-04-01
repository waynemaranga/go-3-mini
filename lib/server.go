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

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func StartServer(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", chatHandler)
	mux.HandleFunc("/prompt", promptHandler)
	mux.HandleFunc("/chats", chatsHandler)

	LogInfo("Server starting on port " + port)
	return http.ListenAndServe(":"+port, mux)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	history := GetChatHistory()
	userMessage := ChatMessage{Role: "user", Content: req.Message}
	history = append(history, userMessage)
	SaveChat(userMessage)

	aiResponse := GetAIResponse(history)
	aiMessage := ChatMessage{Role: "assistant", Content: aiResponse}
	SaveChat(aiMessage)

	json.NewEncoder(w).Encode(map[string]string{"response": aiResponse})
}

func promptHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	messages := []ChatMessage{
		{Role: "user", Content: req.Prompt},
	}

	response := GetAIResponse(messages)
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func chatsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history := GetChatHistory()
	json.NewEncoder(w).Encode(history)
}
