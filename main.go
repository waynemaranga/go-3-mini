package main

import (
	"fmt"
	"go-3-mini/lib"
)

func main() {
	fmt.Println("Starting Chatbot...")
	lib.ConnectDB()
	lib.StartShell()
}
