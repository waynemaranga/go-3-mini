package lib

import "os"

var (
	// MongoURI   = os.Getenv("MONGO_URI") // Set this in your environment
	MongoURI   = "mongodb://localhost:27017"
	DBName     = "chatbot"
	Collection = "messages"
	OpenAIKey  = os.Getenv("OPENAI_API_KEY") // Set this in your environment
)
