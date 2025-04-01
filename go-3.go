package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages            []Message `json:"messages"`
	MaxCompletionTokens int       `json:"max_completion_tokens"`
	Model               string    `json:"model"`
}

type ChatResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func createCompletion(userMessage string, systemMessage string) (*ChatResponse, error) {
	endpoint := os.Getenv("ENDPOINT")
	apiKey := os.Getenv("API_KEY")

	messages := []Message{
		{
			Role:    "system",
			Content: systemMessage,
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}

	chatRequest := ChatRequest{
		Messages:            messages,
		MaxCompletionTokens: 100000,
		Model:               "o3-mini",
	}

	jsonData, err := json.Marshal(chatRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	url := endpoint + "/openai/deployments/o3-mini/chat/completions?api-version=2025-01-01-preview"
	// use TARGET_URI from .env
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: API returned status code %d: %s", resp.StatusCode, string(body))
	}

	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &chatResponse, nil
}

func main() {
	err := godotenv.Load("./go-3-mini.env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	userMessage := "Hello, how are you?"
	systemMessage := "You are a helpful assistant."
	
	completion, err := createCompletion(userMessage, systemMessage)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println(completion)
	fmt.Println("🐬")
}