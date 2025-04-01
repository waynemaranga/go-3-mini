package lib

import (
	"bufio"
	"fmt"
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

		userMessage := ChatMessage{Role: "user", Content: input}
		SaveChat(userMessage)

		history := GetChatHistory()
		aiResponse := GetAIResponse(history)

		botMessage := ChatMessage{Role: "assistant", Content: aiResponse}
		SaveChat(botMessage)

		fmt.Println("Bot:", aiResponse)
	}
}
