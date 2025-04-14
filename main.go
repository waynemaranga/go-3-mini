package main

import (
	"fmt"
	"go-3-mini/lib"
	"os"
)

func main() {
	fmt.Println("⏳ Starting go-3-mini...")
	lib.ConnectDB()
	lib.InitCache() // Initialize the cache

	// Check command-line arguments
	if len(os.Args) > 1 && os.Args[1] == "shell" {
		fmt.Println("💬 Starting interactive shell...")
		lib.StartShell()
		return
	}

	// Default to starting the HTTP server
	fmt.Println("🚀 Starting HTTP server on port 8080...")
	if err := lib.StartServer("8080"); err != nil {
		fmt.Printf("⛔ Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
