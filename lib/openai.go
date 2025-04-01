package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	MaxCompletionTokens int           `json:"max_completion_tokens,omitempty"`
}

func GetAIResponse(history []ChatMessage) string {
	// Load environment variables
	godotenv.Load(".env")
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	targetURI := os.Getenv("AZURE_OPENAI_TARGET_URI")

	// Set up the request
	requestBody, _ := json.Marshal(OpenAIRequest{
		Model:               "o3-mini",
		Messages:            history,
		MaxCompletionTokens: 100000,
	})

	req, _ := http.NewRequest("POST", targetURI, bytes.NewBuffer(requestBody))
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling API:", err)
		return "Error"
	}
	defer resp.Body.Close()

	// Print the full raw response
	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw Response:", string(bodyBytes))

	return "Response logged"
}
