// db/db.go
package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB client wrapper
type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// Connect establishes a connection to MongoDB
func Connect(ctx context.Context, uri string) (*MongoDB, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client: client,
		db:     client.Database("mini-chat"),
	}, nil
}

// GetCollection returns a MongoDB collection
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

// FindChats retrieves chats from the database
func (m *MongoDB) FindChats(ctx context.Context) (*mongo.Cursor, error) {
	return m.GetCollection("chats").Find(ctx, bson.M{})
}

// FindChatByID finds a chat by its ID
func (m *MongoDB) FindChatByID(ctx context.Context, id interface{}) *mongo.SingleResult {
	return m.GetCollection("chats").FindOne(ctx, bson.M{"_id": id})
}

// InsertChat inserts a new chat
func (m *MongoDB) InsertChat(ctx context.Context, chat interface{}) (*mongo.InsertOneResult, error) {
	return m.GetCollection("chats").InsertOne(ctx, chat)
}

// UpdateChat updates an existing chat
func (m *MongoDB) UpdateChat(ctx context.Context, id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return m.GetCollection("chats").UpdateOne(ctx, bson.M{"_id": id}, update)
}

// Close disconnects from MongoDB
func (m *MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
