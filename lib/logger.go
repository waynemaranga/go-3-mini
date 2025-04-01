package lib

import "fmt"

func LogInfo(message string) {
	fmt.Println("[INFO]", message)
}

func LogError(err error) {
	if err != nil {
		fmt.Println("[ERROR]", err)
	}
}
