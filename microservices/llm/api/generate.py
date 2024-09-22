# /api/generate.py
from flask import Blueprint, request, jsonify
from clients.openai_client import generate_text

# Blueprint for the API routes
generate_api = Blueprint('generate_api', __name__)

@generate_api.route('/generate', methods=['POST'])
def generate():
    data = request.json
    if not data or 'prompt' not in data:
        return jsonify({'error': 'No prompt provided'}), 400
    
    prompt = data['prompt']
    result = generate_text(prompt)

    if "Error" in result:
        return jsonify({'error': result}), 500

    return jsonify({'generated_text': result})
