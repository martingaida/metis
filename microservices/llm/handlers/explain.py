from dotenv import load_dotenv
from pydantic import BaseModel
from openai import OpenAI
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
    layers: List[Layer]

class Topic(BaseModel):
    topic: str
    concepts: List[Concept]

class StructuredExplanation(BaseModel):
    topics: List[Topic]
    main_takeaway: str


def generate_structured_explanation(text):
    prompt = f"""Analyze the following text and provide a structured explanation:

        1. Identify the main topics discussed in the text. The number of topics should reflect the content; there may be one dominant topic or multiple important topics.
        2. For each topic, list the key concepts. The number of concepts should be appropriate to fully represent the topic without redundancy.
        3. For each concept, provide a three-layered explanation using the 'What, Why, How' framework:
        - Layer 1: Simple explanation for beginners
        - Layer 2: Detailed explanation with examples
        - Layer 3: Technical explanation for advanced readers
        4. Provide a short, concise summary (1-3 sentences) that captures the most significant takeaway from the entire text.

        Present the response in the following JSON format:

        {{
            "topics": [
                {{
                "topic": "Topic Name",
                "concepts": [
                    {{
                    "concept": "Concept Name",
                    "layers": [
                        {{
                        "layer_1": {{
                            "what": "Simple explanation of what it is",
                            "why": "Simple explanation of why it's important",
                            "how": "Simple explanation of how it works"
                        }}
                        }},
                        {{
                        "layer_2": {{
                            "what": "More detailed explanation with examples",
                            "why": "Detailed explanation of its significance",
                            "how": "More detailed explanation of its mechanics"
                        }}
                        }},
                        {{
                        "layer_3": {{
                            "what": "Technical definition and components",
                            "why": "In-depth explanation of its importance",
                            "how": "Technical explanation of its workings"
                        }}
                        }}
                    ]
                    }}
                ]
                }}
            ],
            "main_takeaway": "A concise summary capturing the most important point from the entire text."
        }}

        Text to analyze:
        {text}
    """

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

def generate_response(text):
    try:
        start_time = time.time()
        result = generate_structured_explanation(text)
        total_time = time.time() - start_time
        print(f'Total function runtime: {total_time:.2f} seconds')
        return {"explanations": result.model_dump()}
    except Exception as e:
        return f"Error: {str(e)}"


# # For testing
# if __name__ == "__main__":
#     sample_text = """Abstract. Rapid advances in Artificial Intelligence (AI) are
#         generating much controversy in society, often without scientific
#         basis. As occurred the development of other emerging technologies,
#         such as the introduction of electricity in the early 20th century, AI
#         causes both fascination and fear. Following the advice of the
#         philosopher R.W. Emerson's advice ‘the knowledge is the antidote
#         to fear’, this paper seeks to contribute to the dissemination of
#         knowledge about AI. To this end, it reflects on the following
#         questions: the origins of AI, its possible future evolution, its ability
#         to show feelings, the associated threats and dangers, and the
#         concept of AI singularity
#         Conclusions
#         AI is emerging as the main focus of the sixth industrial revolution and
#         a key catalyst in the emergence of a new world order. This paper reflects
#         on several crucial aspects, such as the origins of AI, its possible future
#         evolution, its ability to simulate feelings, the associated threats and
#         dangers, and the concept of AI singularity. The main conclusions are
#         presented below.
#         Defining AI as the discipline that makes machines perform tasks
#         that, if done by humans, would require intelligence, we can place its
#         origins in the Early Modern Era, when the first calculating machines
#         were created.
#         In terms of future developments, some of the challenges would be to
#         develop machines capable of emulating human mental capacities, such
#         as reasoning, comprehension, imagination, perception, recognition,
#         creativity and emotions. Although we are still far from achieving these
#         goals, very significant partial progress has been made. On the other
#         hand, while AI systems can simulate emotions in a useful way in certain
#         contexts, they are not capable of experiencing real feelings.
#         Among the current, non-speculative problems related to AI, several
#         major drawbacks stand out. One of the most discussed is the
#         destruction of jobs, requiring the development of new training and
#         adaptation strategies. In addition, there exists a propagandistic misuse
#         of the term ‘AI’, attributing it to systems that actually do not comply
#         with its characteristics. Another relevant problem is the global
#         monitoring and control of data, which makes it possible to extract 
#         Five questions and answers about artificial intelligence 15
#         confidential information and even alter the balance of power at the
#         global level, given that information is a strategic resource. This is
#         compounded by the creation and dissemination of fake news with the
#         appearance of authenticity, impersonation, and the increasing difficulty
#         of tracing its origin and assigning responsibility in these processes.
#         Significant progress has been made towards the creation of an
#         artificial general intelligence (AGI), i.e. an AI with a flexible intellect
#         comparable to humans, thanks to the development of OpenAI's GPT-4
#         language model, which allows machines to converse with each other as
#         if they were human beings [20]. This advance represents an approach
#         to one of the biggest potential risks: the feared singularity of AI.
#         Speculatively, it is believed that at this point a superintelligence capable
#         of monitoring and controlling all aspects of reality could be achieved.
#         Despite the above dangers, it is essential to remember that
#         everything a computer system, and especially an AI programme, does
#         is the result of what its designers intended. A machine, on its own,
#         cannot execute any operation that has not been planned and
#         anticipated by human beings. In fact, intelligence and ingenuity reside
#         not in the algorithms and machines, but in the people, who conceive
#         and develop them.
#         AI does not possess intelligence in the human sense; it is limited to
#         learning and making deductions. It lacks creativity and is not able to
#         hypothesise, speculate, make discoveries on its own initiative, or
#         automatically apply its abilities to different areas, as we humans do.
#         Computers can reason, but they do not think; if they are not given initial
#         data or information, they do not know what to do on their own. To reach
#         a human-like level of intelligence, machines would have to be
#         autonomously intuitive and creative.
#         AI systems must be transparent, subject to human oversight, and
#         assessable and certifiable by external authorities. It is essential to
#         ensure that the data used to train these systems is free of bias and that
#         fundamental rights are always respected. The key challenge is to
#         achieve developments that truly drive progress, equality and prosperity
#         for all, not just for the few. To this end, it is crucial to have adequate
#         legislation and, most difficult of all, to establish effective monitoring
#         systems to ensure compliance. Fortunately, various entities such as
#         UNESCO, the European Union and several states are already working
#         in this direction."""

#     response = generate_response(sample_text)
#     print(json.dumps(response, indent=2))