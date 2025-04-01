# `go-3-mini`
## Chat Application Using Azure OpenAI's o3-mini Model

This project is a Go-based chat application using Azure's OpenAI Service to interact with the `o3-mini` reasoning model. The application facilitates conversations by sending user inputs to the model and displaying the generated responses.

## Project Structure

The project is organized as follows:

- `main.go`: The entry point of the application, managing user interactions and displaying chat responses.
- `lib/`: Contains the core functionality of the application.
  - `lib/ai_client.go`: Handles communication with the Azure OpenAI Service.
  - `lib/models.go`: Defines data structures for chat messages and API requests/responses.

## Features of the o3-mini Model

The `o3-mini` model offers:

- **Enhanced Reasoning**: Improved problem-solving capabilities, particularly in STEM fields such as science, mathematics, and coding.
- **Cost Efficiency**: More affordable per-token pricing compared to earlier models like `o1-mini`, making it suitable for large-scale applications.
- **Low Latency**: Faster response times, enhancing user experience in interactive applications.

For more details, refer to the [Azure OpenAI Service announcement](https://azure.microsoft.com/en-us/blog/announcing-the-availability-of-the-o3-mini-reasoning-model-in-microsoft-azure-openai-service/).

## Prerequisites

Before running the application, ensure you have:

- A Microsoft Azure account with access to the Azure OpenAI Service. <https://azure.microsoft.com/en-gb>
- An active deployment of the `o3-mini` model within your Azure subscription. <https://ai.azure.com>
- Go programming language installed on your system. <https://go.dev/doc/install>
- Necessary environment variables configured for API access. See [`.env.example`](./.env.example) for reference.

## Setup Instructions

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/waynemaranga/go-3-mini.git
   cd go-3-mini
   ```

2. **Install Dependencies**:

   ```bash
   go mod tidy
   ```

3. **Configure Environment Variables**:

   Duplicate `.env.example` as `.env` and fill in the required values from Azure AI Foundry:

   ```env
   AZURE_OPENAI_API_KEY=your_api_key_here
   AZURE_OPENAI_TARGET_URI=https://your_endpoint_here/openai/deployments/o3-mini/chat/completions?api-version=2025-01-01-preview
   ```

4. **Run the Application**:

   ```sh
   go run main.go
   ```

## How It Works

1. **User Input**: The application prompts the user to enter a message.
2. **API Request**: The input is sent to the `o3-mini` model via the Azure OpenAI Service.
3. **Response Handling**: The model generates a response, which is displayed to the user.

## Integration with Azure OpenAI Service

The application communicates with the Azure OpenAI Service using HTTP requests. The key integration points are:

- **API Endpoint**: Defined by the `AZURE_OPENAI_TARGET_URI` environment variable.
- **Authentication**: Managed via the `AZURE_OPENAI_API_KEY` environment variable.
- **Request Structure**: The application constructs a JSON payload containing the model name, input messages, and other parameters, which is then sent to the API endpoint.

## Additional Resources

For more information on the `o3-mini` model and Azure OpenAI Service, refer to:

- [Azure OpenAI Service Models](https://learn.microsoft.com/en-us/azure/ai-services/openai/concepts/models)
- [Azure AI Foundry Model Catalog](https://ai.azure.com/explore/models/o3-mini/version/2025-01-31/registry/azure-openai)

## License

This project is licensed under the MIT License. 