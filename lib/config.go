package lib

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MongoURI             string
	DBName               = "go_3_mini"
	Collection           = "chats"
	AzureOpenAIAPIKey    string
	AzureOpenAITargetURI string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("⛔ Error loading .env file")
	}

	MongoURI = os.Getenv("MONGODB_URI")
	AzureOpenAIAPIKey = os.Getenv("AZURE_OPENAI_API_KEY")
	AzureOpenAITargetURI = os.Getenv("AZURE_OPENAI_TARGET_URI")
}
