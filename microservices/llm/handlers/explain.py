from dotenv import load_dotenv
from openai import OpenAI
from pathlib import Path
import os

# Get the parent directory path
microservices_root_dir = Path(__file__).resolve().parent.parent.parent

# Load the .env file from the parent directory
load_dotenv(dotenv_path=microservices_root_dir / ".env")

# Set up Open AI client
client = OpenAI()


def synthesize_major_topics(text):
    print("Synthesizing major topics...")
    prompt = f"Analyze the following text and identify the overarching major topic in clear, simple language. If there is one clear, dominant topic, provide only that single topic. If there are multiple equally important topics, list them. Keep the explanation concise and aimed at a general audience. Text:\n\n{text}\n\n"
    response = client.chat.completions.create(
        model="gpt-4",
        messages=[
            {"role": "system", "content": "You are an expert at synthesizing complex information into easily understandable topics."},
            {"role": "user", "content": prompt}
        ]
    )
    topics = response.choices[0].message.content.split('\n')
    print(f'{len(topics)} topics found')
    return [topic.strip() for topic in topics if topic.strip()]


def synthesize_major_concepts(topic):
    print("Synthesizing major concepts...")
    prompt = f"For the topic '{topic}', break it down into its simplest major concepts. List each major concept, using clear and accessible language. Ensure that these concepts are understandable to someone unfamiliar with the field."
    response = client.chat.completions.create(
        model="gpt-4",
        messages=[
            {"role": "system", "content": "You are an expert at breaking down topics into core concepts for non-expert audiences."},
            {"role": "user", "content": prompt}
        ]
    )
    concepts = response.choices[0].message.content.split('\n')
    print(f'{len(concepts)} concepts found')
    return [concept.strip() for concept in concepts if concept.strip()]


def generate_layered_explanation(concept):
    print("Generating layers of explanation...")
    layers = []

    for level in range(1, 4):
        if level == 1:
            prompt = f"Explain what {concept} is, why it's important, and how it works in the simplest terms possible, as if explaining to a complete beginner or non-scientific audience."
        elif level == 2:
            prompt = f"Provide a more detailed explanation of {concept}. Include examples and practical applications to help a general audience understand why it is significant and how it works."
        else:
            prompt = f"Give a thorough, technical explanation of {concept}, including its components, importance, and real-world applications. This explanation should be suitable for someone with more advanced knowledge of the subject."

        response = client.chat.completions.create(
            model="gpt-4",
            messages=[
                {"role": "system", "content": "You are a knowledgeable tutor explaining complex topics in simple and technical layers."},
                {"role": "user", "content": prompt}
            ]
        )
        layers.append(f"{response.choices[0].message.content}")
    
    return layers


def process_text(text):
    print("Processing text...")
    """ Returns:
        [
            {
                "topic": "major_topic_1",
                "concepts": [
                    {
                        "concept": "concept_1",
                        "layers": [
                            {"layer_1": "layer_1_explanation"},
                            {"layer_2": "layer_2_explanation"},
                            {"layer_3": "layer_3_explanation"}
                        ]
                    },
                    # More concepts...
                ]
            },
            # More topics...
        ]
    """
    result = []
    
    # Synthesize major topics
    major_topics = synthesize_major_topics(text)
    
    # Process each major topic
    for topic in major_topics:
        topic_data = {
            "topic": topic,
            "concepts": []
        }

        # Synthesize major concepts for each topic
        major_concepts = synthesize_major_concepts(topic)
        
        # Generate layered explanations for each concept
        for concept in major_concepts:
            concept_data = {
                "concept": concept,
                "layers": []
            }            
            
            layers = generate_layered_explanation(concept)

            for i, layer in enumerate(layers, 1):
                concept_data["layers"].append({f"layer_{i}": layer})
            
            topic_data["concepts"].append(concept_data)
        
        result.append(topic_data)
    
    return result


def generate_response(text):
    try:
        result = process_text(text)
        return result
    except Exception as e:
        return f"Error: {str(e)}"

# if __name__ == "__main__":
#     sample_text = """
#         Solid-state spin systems hold great promise for quantum information processing and the construction of quantum networks. However, the considerable inhomogeneity of spins in solids poses
#         a significant challenge to the scaling of solid-state quantum systems. A practical protocol to individually control and entangle spins remains elusive. To this end, we propose a hybrid spin-phonon
#         architecture based on spin-embedded SiC optomechanical crystal (OMC) cavities, which integrates
#         photonic and phononic channels allowing for interactions between multiple spins. With a Ramanfacilitated process, the OMC cavities support coupling between the spin and the zero-point motion
#         of the OMC cavity mode reaching 0.57 MHz, facilitating phonon preparation and spin Rabi swap
#         processes. Based on this, we develop a spin-phonon interface that achieves a two-qubit controlledZ gate with a simulated fidelity of 96.80% and efficiently generates highly entangled Dicke states
#         with over 99% fidelity, by engineering the strongly coupled spin-phonon dark state which is robust
#         against loss from excited state relaxation as well as spectral inhomogeneity of the defect centers.
#         This provides a hybrid platform for exploring spin entanglement with potential scalability and full
#         connectivity in addition to an optical link, and offers a pathway to investigate quantum acoustics
#         in solid-state systems.
#     """
#     result = generate_response(sample_text)
#     print(result)