import os
from openai import AzureOpenAI
import dotenv
from typing import Optional
from openai.types.chat.chat_completion import ChatCompletion

dotenv.load_dotenv("./go-3-mini.env")
ENDPOINT: Optional[str] = os.getenv("ENDPOINT")
API_KEY: Optional[str] = os.getenv("API_KEY")

client = AzureOpenAI(
    api_version="2025-01-01-preview",
    azure_endpoint=ENDPOINT,
    api_key=API_KEY,
    )


def create_completion(
        user_message: str,
        system_message: str = "You are a helpful assistant."
        ) -> ChatCompletion:
    
    response: ChatCompletion = client.chat.completions.create(
    messages=[
        {"role": "system","content": system_message,},
        {"role": "user","content": user_message,}
        ],
        max_completion_tokens=100000,
        model="o3-mini"
        )

    return response


if __name__ == "__main__":
    user_message: str = "Hello, how are you?"
    completion: ChatCompletion = create_completion(user_message)
    # print(completion.choices[0].message.content)
    print(completion)
    print("🐬")