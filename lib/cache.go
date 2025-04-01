// lib/cache.go
package lib

import (
	"os"
	"strconv"
	"sync"
)

// In-memory cache for chat history
var (
	chatHistory  []ChatMessage
	mutex        sync.RWMutex
	maxCacheSize int = 1000 // Default max cache size
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

func init() {
	if sizeStr := os.Getenv("CACHE_SIZE"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil {
			maxCacheSize = size
		}
	}

	// Can also check CACHE_ENABLED if you want to make caching optional
	// cacheEnabled := os.Getenv("CACHE_ENABLED") == "true"
}
