package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model               string        `json:"model"`
	Messages            []ChatMessage `json:"messages"`
	MaxCompletionTokens int           `json:"max_tokens,omitempty"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func GetAIResponse(history []ChatMessage) string {
	// Load environment variables
	godotenv.Load(".env")
	// endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	targetURI := os.Getenv("AZURE_OPENAI_TARGET_URI")

	// Create request with system message if not already in history
	messages := history
	hasSystemMessage := false
	for _, msg := range messages {
		if msg.Role == "system" {
			hasSystemMessage = true
			break
		}
	}

	if !hasSystemMessage {
		systemMessage := ChatMessage{
			Role:    "system",
			Content: "You are a helpful assistant.",
		}
		messages = append([]ChatMessage{systemMessage}, messages...)
	}

	// Set up the request
	// url := endpoint + "/openai/deployments/o3-mini/chat/completions?api-version=2025-01-01-preview"
	url := targetURI

	requestBody, _ := json.Marshal(OpenAIRequest{
		Model:               "o3-mini",
		Messages:            messages,
		MaxCompletionTokens: 100000,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling API:", err)
		return "Error"
	}
	defer resp.Body.Close()

	var openAIResp OpenAIResponse
	json.NewDecoder(resp.Body).Decode(&openAIResp)
	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content
	}
	return "No response"
}
