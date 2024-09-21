# /app.py
from flask import Flask
from api.generate import generate_api
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

app = Flask(__name__)

# Register the blueprint
app.register_blueprint(generate_api)

if __name__ == "__main__":
    port = os.getenv("PORT", 5000)
    app.run(host="0.0.0.0", port=port, debug=True)