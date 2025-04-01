# `go-3-mini`

A lightweight chat application for OpenAI's o3-mini from Azure AI with MongoDB storage.

## Features

- Interactive chat shell
- HTTP server with REST API
- Conversation history stored in MongoDB
- Azure OpenAI integration for AI responses

## Modules

### Core Modules
- **config.go**: Loads environment variables (MongoDB URI, OpenAI API keys)
- **logger.go**: Simple logging utilities
- **mongo.go**: MongoDB connection and chat history management
- **openai.go**: Azure OpenAI API integration
- **server.go**: HTTP server with REST endpoints
- **shell.go**: Interactive command-line interface

## Quick Start

### Prerequisites
- Go 1.20+
- MongoDB (local or remote)
- Azure OpenAI API access

### Setup
1. Clone the repository
2. Create `.env` file:
   ```env
   MONGODB_URI=mongodb://localhost:27017
   AZURE_OPENAI_API_KEY=your_api_key
   AZURE_OPENAI_TARGET_URI=https://your-endpoint.openai.azure.com/...
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application
Start with either interactive shell or HTTP server:
```bash
go run main.go
```

Choose mode when prompted:
```bash
1️⃣. Start HTTP server
2️⃣. Start interactive shell
```

### HTTP Server Endpoints
- `POST /chat` - Chat with history (expects `{"message": "your text"}`)
- `POST /prompt` - Single prompt (expects `{"prompt": "your text"}`)
- `GET /chats` - Get all chat history

### Interactive Shell
Type messages directly and get AI responses. Type "exit" to quit.

## Configuration
Edit `.env.example` file into `.env` to configure:
- MongoDB connection string
- Azure OpenAI credentials
- Database name (default: `go_3_mini`)

## License
MIT. See [LICENSE.md](./LICENSE.d) for more information.


