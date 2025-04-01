import os
import json
import http.client
from urllib.parse import urlparse
from typing import Optional, Any, TypedDict, Union

from dotenv import load_dotenv
load_dotenv()

ENDPOINT: Optional[str] = os.getenv("AZURE_OPENAI_TARGET_URI")
API_KEY: Optional[str] = os.getenv("AZURE_OPENAI_API_KEY")

class Message(TypedDict):
    role: str
    content: str

class Choice(TypedDict):
    message: Message

class ChatCompletion(TypedDict):
    choices: list[Choice]

def create_completion(
    user_message: str,
    system_message: str = "You are a helpful assistant."
    ) -> ChatCompletion:

    if ENDPOINT is None or API_KEY is None:
        raise ValueError("AZURE_OPENAI_ENDPOINT and AZURE_OPENAI_API_KEY must be set as environment variables.")

    payload: dict[str, Union[str, list[Message], int]] = {
        "model": "o3-mini",
        "messages": [
            {"role": "system", "content": system_message},
            {"role": "user", "content": user_message}
        ],
        "max_tokens": 100000
    }

    parsed_url = urlparse(ENDPOINT)
    conn = http.client.HTTPSConnection(parsed_url.netloc)

    payload_json = json.dumps(payload)

    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {API_KEY}'
    }

    conn.request("POST", parsed_url.path, body=payload_json, headers=headers)
    response = conn.getresponse()
    response_data = response.read().decode('utf-8')
    conn.close()

    return json.loads(response_data)

if __name__ == "__main__":
    user_message: str = input("Enter your message: ")
    completion: ChatCompletion = create_completion(user_message)
    # print(f"\nResponse: {completion['choices'][0]['message']['content']}")
    print(f"\nResponse: {completion}")
    print("🐬")
