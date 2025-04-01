// clients/openai.go
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OpenAIBaseURL = "https://api.openai.com/v1"
)

// Message represents a chat message for the OpenAI API
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the request structure for the OpenAI chat API
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse is the response structure from OpenAI chat API
type ChatResponse struct {
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

// Client is the OpenAI API client
type Client struct {
	APIKey string
	Model  string
	client *http.Client
}

// NewClient creates a new OpenAI client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		Model:  "o3-mini",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateChatCompletion sends a chat completion request to OpenAI
func (c *Client) CreateChatCompletion(messages []Message) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("OpenAI API key not provided")
	}

	reqBody := ChatRequest{
		Model:    c.Model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest("POST", OpenAIBaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}
