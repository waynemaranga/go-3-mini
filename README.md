# `go-3-mini`

<!-- Badges -->

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Azure](https://img.shields.io/badge/microsoft%20azure-0089D6?style=for-the-badge&logo=microsoft-azure&logoColor=white)
![Angular](https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
![Typescript](https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white)

A lightweight chat application for OpenAI's **o3-mini** from Azure AI with MongoDB storage. The frontend is built with Angular and the backend is written in Go.

## 📚 Table of Contents

- [`go-3-mini`](#go-3-mini)
  - [📚 Table of Contents](#-table-of-contents)
  - [☑️ Features](#️-features)
  - [📦 Modules](#-modules)
    - [📦 Core Modules](#-core-modules)
  - [🏁 Quick Start](#-quick-start)
    - [🧰 Prerequisites](#-prerequisites)
    - [Azure OpenAI API](#azure-openai-api)
    - [🛠️ Setup](#️-setup)
    - [⏳ Running the Application](#-running-the-application)
    - [HTTP Server Endpoints](#http-server-endpoints)
    - [Interactive Shell](#interactive-shell)
  - [Configuration](#configuration)
  - [🐋 Docker](#-docker)
    - [🧰 Prerequisites](#-prerequisites-1)
    - [Quick Start](#quick-start)
    - [Features](#features)
    - [Configuration](#configuration-1)
    - [Commands](#commands)
    - [Production Notes](#production-notes)
    - [Troubleshooting](#troubleshooting)
  - [License](#license)

## ☑️ Features

- Interactive chat shell
- HTTP server with REST API
- Conversation history stored in MongoDB
- Azure OpenAI integration for AI responses
- Web interface

## 📦 Modules

### 📦 Core Modules

- **`config.go`**: Loads environment variables (MongoDB URI, OpenAI API keys)
- **`logger.go`**: Simple logging utilities
- **`mongo.go`**: MongoDB connection and chat history management
- **`openai.go`**: Azure OpenAI API integration
- **`server.go`**: HTTP server with REST endpoints
- **`shell.go`**: Interactive command-line interface
- **`ui`**: Angular frontend for web interface

## 🏁 Quick Start

### 🧰 Prerequisites

- Go 1.23+
- MongoDB (local or remote)
- Azure OpenAI API access
- Node.js 20+

### Azure OpenAI API

1. Create an Azure account
   - You can use the free tier for OpenAI or
   - **You can use Microsoft for Startups to get $1000 in Azure credits** <https://www.microsoft.com/en-us/startups>
2. Create an OpenAI resource
   - In Azure you can use Azure AI Foundry at <https://ai.azure.com>
3. Get the API key and endpoint URL


### 🛠️ Setup

1. Clone the repository
2. Create `.env` file from `.env.example` and configure:
   ```env
   MONGODB_URI=mongodb://localhost:27017
   AZURE_OPENAI_API_KEY=your_api_key
   AZURE_OPENAI_TARGET_URI=https://your-endpoint.openai.azure.com/...
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### ⏳ Running the Application

0. (Locally) Start MongoDB server:

```bash
mkdir -p .mongo/
mongod --dbpath .mongo/
```

1. Start app with either interactive shell or HTTP server:

```bash
go run main.go
```

2. Choose mode when prompted:

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

## 🐋 Docker

### 🧰 Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- `.env` file configured (see [Configuration](#configuration))

### Quick Start

1. Build and launch containers:
   ```bash
   docker-compose up --build
   ```
2. Access the application:
   - CLI: Connect to the container shell
   - Web: http://localhost:8080

### Features

- 🐳 Full containerization of app and MongoDB
- 🔐 Secure MongoDB initialization with:
  - Root user authentication
  - Dedicated application user
  - Automatic database/collection creation
- 🩺 Health checks for both services
- ♻️ Automatic restarts on failure

### Configuration

Configure via `.env` file or environment variables:

| Variable             | Default   | Description                   |
| -------------------- | --------- | ----------------------------- |
| `APP_PORT`           | 8080      | Application port              |
| `MONGO_PORT`         | 27017     | MongoDB exposed port          |
| `MONGO_USER`         | root      | MongoDB root username         |
| `MONGO_PASSWORD`     | example   | MongoDB root password         |
| `MONGO_APP_USER`     | appuser   | Application database user     |
| `MONGO_APP_PASSWORD` | apppass   | Application database password |
| `DB_NAME`            | go_3_mini | Database name                 |
| `COLLECTION`         | chats     | Default collection            |

### Commands

**Start in detached mode:**

```bash
docker-compose up -d
```

**View logs:**

```bash
docker-compose logs -f
```

**Stop services:**

```bash
docker-compose down
```

**Remove volumes (caution - deletes data):**

```bash
docker-compose down -v
```

**Rebuild containers:**

```bash
docker-compose up --force-recreate --build
```

### Production Notes

1. **Secrets Management**:

   - Use Docker secrets or mounted secret files for credentials
   - Never commit `.env` files to version control

2. **Persistent Storage**:

   - MongoDB data persists via named volume (`mongo-data`)
   - Backup volume regularly:
     ```bash
     docker run --rm -v mongo-data:/source -v $(pwd):/backup alpine \
       tar czf /backup/mongo-backup-$(date +%Y%m%d).tar.gz -C /source .
     ```

3. **Scaling**:
   ```bash
   docker-compose up -d --scale app=3  # Multiple app instances
   ```

### Troubleshooting

**Common Issues:**

- **Connection failures**: Verify MongoDB health checks complete first
- **Permission errors**: Ensure `mongo-init.sh` is executable
- **Missing environment variables**: Check `.env` file exists

**Debugging:**

```bash
# Inspect MongoDB initialization:
docker exec -it mongodb cat /docker-entrypoint-initdb.d/mongo-init.sh

# Enter MongoDB shell:
docker exec -it mongodb mongosh -u $MONGO_USER -p $MONGO_PASSWORD
```

## License

MIT. See [LICENSE.md](./LICENSE.d) for more information.
