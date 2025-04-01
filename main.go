package main

import (
	"fmt"
	"go-3-mini/lib"
)

func main() {
	fmt.Println("⏳ Starting go-3-mini...")
	lib.ConnectDB()
	lib.StartShell()
}
