// lib/cache.go
package lib

import (
	"sync"
)

// In-memory cache for chat history
var (
	chatHistory []ChatMessage
	mutex       sync.RWMutex
)

// InitCache loads initial chat history from DB
func InitCache() {
	mutex.Lock()
	defer mutex.Unlock()

	chatHistory = GetChatHistoryFromDB()
	LogInfo("✅ Loaded chat history from database")
}

// AddChat adds a message to in-memory cache and saves to DB
func AddChat(chat ChatMessage) {
	mutex.Lock()
	defer mutex.Unlock()

	chatHistory = append(chatHistory, chat)
	go SaveChatToDB(chat) // Asynchronously save to DB
}

// GetChatHistoryFromCache returns a copy of the in-memory chat history
func GetChatHistoryFromCache() []ChatMessage {
	mutex.RLock()
	defer mutex.RUnlock()

	// Create a copy to prevent race conditions
	historyCopy := make([]ChatMessage, len(chatHistory))
	copy(historyCopy, chatHistory)
	return historyCopy
}
