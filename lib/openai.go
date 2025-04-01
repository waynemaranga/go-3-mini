package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAIRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func GetAIResponse(history []ChatMessage) string {
	url := "https://api.openai.com/v1/chat/completions"
	requestBody, _ := json.Marshal(OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: history,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "Bearer "+OpenAIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling OpenAI:", err)
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
