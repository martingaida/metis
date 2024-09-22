from dotenv import load_dotenv
from openai import OpenAI
from pathlib import Path
import os


# Get the parent directory path
parent_dir = Path(__file__).resolve().parent.parent

# Load the .env file from the parent directory
load_dotenv(dotenv_path=parent_dir / '.env')


def generate_response():
    # Set up Open AI client
    client = OpenAI()

    assistant = client.beta.assistants.create(
        name="Math Tutor",
        instructions="You are a personal math tutor. Write and run code to answer math questions.",
        tools=[{"type": "code_interpreter"}],
        model="gpt-4o",
    )

    thread = client.beta.threads.create()

    message = client.beta.threads.messages.create(
        thread_id=thread.id,
        role="user",
        content="I need to solve the equation `3x + 11 = 14`. Can you help me?"
    )

    run = client.beta.threads.runs.create_and_poll(
        thread_id=thread.id,
        assistant_id=assistant.id,
        instructions="Please address the user as Jane Doe. The user has a premium account."
    )

    if run.status == 'completed': 
        messages = client.beta.threads.messages.list(
            thread_id=thread.id
        )

        # Destructure and print messages
        for message in messages:
            role = message.role
            content = message.content[0].text.value
            print(f"Role: {role}")
            print(f"Content: {content}")
            print("---")
    else:
        print(run.status)

generate_response()