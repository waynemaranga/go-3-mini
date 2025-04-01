package models

type ChatMessage struct {
	Role    string `json:"role" bson:"role"`
	Content string `json:"content" bson:"content"`
}
