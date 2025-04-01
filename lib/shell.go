package lib

import (
	"bufio"
	"fmt"
	"go-3-mini/ai"
	"go-3-mini/db"
	"go-3-mini/models"
	"os"
	"strings"
)

func StartShell() {
	fmt.Println("Chatbot Shell - Type 'exit' to quit")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("You: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		userMessage := models.ChatMessage{Role: "user", Content: input}
		db.SaveChat(userMessage)

		history := db.GetChatHistory()
		aiResponse := ai.GetAIResponse(history)

		botMessage := models.ChatMessage{Role: "assistant", Content: aiResponse}
		db.SaveChat(botMessage)

		fmt.Println("Bot:", aiResponse)
	}
}
