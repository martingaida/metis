from dotenv import load_dotenv
from pydantic import BaseModel
from openai import OpenAI, OpenAIError
from pathlib import Path
from typing import List
import time
import json
import os

# Load the .env file from the parent directory
microservices_root_dir = Path(__file__).resolve().parent.parent.parent
load_dotenv(dotenv_path=microservices_root_dir / ".env")

# Set up Open AI client
client = OpenAI()

class Layer(BaseModel):
    what: str
    why: str
    how: str

class Concept(BaseModel):
    concept: str
    layer: Layer

class Topic(BaseModel):
    topic: str
    concepts: List[Concept]

class StructuredExplanation(BaseModel):
    topics: List[Topic]
    main_takeaway: str


def generate_structured_explanation(text, level):
    """
    Generate a structured explanation of the given text using OpenAI's GPT model.

    This function analyzes the input text and produces an explanation
    of the main topics and concepts within it,
    tailored to the provided reading level based on Flesch-Kincaid Score,
    using "What, Why, How" framework.

    Args:
        text (str): The input text to be analyzed and explained.

    Returns:
        StructuredExplanation: A structured object containing the explanation,
        including topics, concepts, and corresponding explanation layers.

    Raises:
        OpenAIError: If there's an issue with the OpenAI API call.
        JSONDecodeError: If the API response cannot be parsed as JSON.
    """

    prompt = f"""Analyze the following text and provide a structured explanation at the {level} reading level:

    1. Identify the main topics discussed in the text. The number of topics should reflect the content; there may be one dominant topic or multiple important topics.
    2. For each topic, list the key concepts. The number of concepts should be appropriate to fully represent the topic without redundancy.
    3. For each concept, provide a single-layered explanation using the 'What, Why, How' framework, tailored to the {level} reading level:
    
    - What: Explain what the concept is, adjusted to the {level} reading level of understanding.
    - Why: Describe why the concept is important or significant, with complexity suitable for the {level} reading level.
    - How: Explain how the concept works or is applied, with detail appropriate for the {level} reading level.

    Adjust the complexity and depth of each explanation based on the specified reading level using Flesch-Kincaid Scale:
        - K3: kindergarten, ages 5-8, example book "Hooray for Fish!"
        - K6: elementary, ages 8-11, example book "The Gruffalo"
        - K9: middle school, ages 11-14, example book "Harry Potter"
        - K12: high school, ages 14-17, example book "Jurassic Park"
        - College: college, ages 17-20, example book "A Brief History of Time"
        - Graduate: graduate, ages 20+, example book "academic papers"

    4. Provide a short summary that captures the most significant takeaway from the entire text, adjusted to the {level} reading level of understanding.

    Present the response in the following JSON format:

    {{
        "topics": [
            {{
                "topic": "Topic Name",
                "concepts": [
                    {{
                        "concept": "Concept Name",
                        "layer": {{
                            "what": "Explanation of what it is, tailored to the {level} level.",
                            "why": "Explanation of why it's important, tailored to the {level} level.",
                            "how": "Explanation of how it works, tailored to the {level} level."
                        }}
                    }}
                ]
            }}
        ],
        "main_takeaway": "A short summary capturing the most important point from the entire text, adjusted to the {level} level."
    }}

    Text to analyze:
    {text}
    """

    try:
        completion = client.beta.chat.completions.parse(
            model="gpt-4o-2024-08-06",
            messages=[
                {"role": "system", "content": "You are an expert at analyzing and explaining complex topics in a structured, multi-layered format. Identify all relevant topics and concepts without a predetermined number."},
                {"role": "user", "content": prompt}
            ],
            response_format=StructuredExplanation,
        )
        
        response_json = json.loads(completion.choices[0].message.content)
        print("Raw API response:", completion.choices[0].message.content)
        print("Parsed response:", response_json)
        result = StructuredExplanation(**response_json)
        print(f'Result: {result}')
        return result
    except OpenAIError as e:
        print(f"OpenAI API error: {e}")
        raise
    except json.JSONDecodeError as e:
        print(f"JSON parsing error: {e}")
        raise
    except (IndexError, ValueError) as e:
        print(f"Unexpected response structure: {e}")
        raise

def generate_response(text, level):
    try:
        start_time = time.time()
        result = generate_structured_explanation(text, level)
        total_time = time.time() - start_time
        print(f'Total function runtime: {total_time:.2f} seconds')
        return {"explanations": result.model_dump()}
    except Exception as e:
        return f"Error: {str(e)}"