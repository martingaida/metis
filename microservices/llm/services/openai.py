# /services/openai_service.py
import openai
import os

# Load the OpenAI API key from environment variables
openai.api_key = os.getenv("OPENAI_API_KEY")

def generate_text(prompt):
    try:
        response = openai.Completion.create(
            engine="text-davinci-003",  # or another model like GPT-4
            prompt=prompt,
            max_tokens=100
        )
        return response.choices[0].text.strip()
    except Exception as e:
        return f"Error: {str(e)}"