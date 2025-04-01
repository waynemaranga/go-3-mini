// lib/shell.go
package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartShell() {
	fmt.Println("🐻 go-3-mini - Type 'exit' to quit")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("🐨 You: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		userMessage := ChatMessage{Role: "user", Content: input}
		AddChat(userMessage)

		history := GetChatHistoryFromCache()
		aiResponse := GetAIResponse(history)

		botMessage := ChatMessage{Role: "assistant", Content: aiResponse}
		AddChat(botMessage)

		fmt.Println("🐻 o3-mini:", aiResponse)
	}
}
