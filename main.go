package main

import (
	"bufio"
	"fmt"
	"go-3-mini/lib"
	"os"
	"strings"
)

func main() {
	fmt.Println("⏳ Starting go-3-mini...")
	lib.ConnectDB()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("☑️ Choose an option:")
	fmt.Println("1️⃣. Start HTTP server")
	fmt.Println("2️⃣. Start interactive shell")
	fmt.Print("Enter your choice (1 or 2): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Println("🚀 Starting HTTP server on port 5467...")
		lib.StartServer("5467")
	case "2":
		lib.StartShell()
	default:
		fmt.Println("⛔ Invalid choice. Please run again and enter either 1 or 2.")
		os.Exit(1)
	}
}
