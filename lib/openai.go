package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type OpenAIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func GetAIResponse(history []ChatMessage) string {
	// Set up the request
	requestBody, _ := json.Marshal(OpenAIRequest{
		Model:               "o3-mini",
		Messages:            history,
		MaxCompletionTokens: 100000,
	})

	req, _ := http.NewRequest("POST", AzureOpenAITargetURI, bytes.NewBuffer(requestBody))
	req.Header.Set("api-key", AzureOpenAIAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("⛔ Error calling API:", err)
		return "⛔ Error calling API: " + err.Error()
	}
	defer resp.Body.Close()

	// Read and parse response
	bodyBytes, _ := io.ReadAll(resp.Body)
	// fmt.Println("Raw Response:", string(bodyBytes)) // TODO: log this

	var openAIResp OpenAIResponse
	json.Unmarshal(bodyBytes, &openAIResp)

	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content
	}
	return "⛔ Error: No response"

}
