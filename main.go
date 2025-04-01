package main

import (
	"fmt"
	"go-3-mini/db"
	"go-3-mini/shell"
)

func main() {
	fmt.Println("Starting Chatbot...")
	db.ConnectDB()
	shell.StartShell()
}
