<p align="center">
  <img src="frontend/src/assets/images/plex-logo.svg" alt="PleX Logo" width="220">
</p>

# PleX: Complex Ideas Made Simple

PleX (short for "Please Explain") is an innovative application designed to explain complex subjects by breaking them down into digestible components. The system works by:

1. Dividing the main subject into major topics
2. Breaking down each topic into key concepts
3. Explaining each concept in three layers of increasing detail:
   - Layer 1: High-level overview
   - Layer 2: Detailed explanation
   - Layer 3: Technical deep-dive
4. Each layer of explanation follows the "What, Why, How" framework to ensure comprehensive understanding:
   - What: Definition and basic description of the concept
   - Why: Importance and relevance of the concept
   - How: Practical application or functioning of the concept

This structured approach allows users to grasp complex information at their preferred depth of understanding, making it suitable for learners at all levels.

## Technical Stack

- Frontend: Angular, TypeScript
- Backend: Go
- Microservices: Python
- Deployment: AWS Lambda, S3

### Architecture Decisions

1. Direct Lambda Invocation:
   - We've chosen to make direct calls to AWS Lambda functions instead of using API Gateway. This decision was made to reduce latency and costs associated with API Gateway, as our application doesn't require the additional features provided by API Gateway such as request throttling or API key management. This approach was especially beneficial before we managed to reduce latency by 90% through optimizing LLM calls and prompt handling.

2. Microservices Architecture:
   - The backend is implemented in Go for its performance benefits and ease of deployment as a Lambda function. However, the OpenAI client libraries are primarily available for Node.js and Python. To leverage these libraries effectively, we've implemented the LLM (Language Model) interaction as a separate microservice in Python. This architecture allows us to combine the performance benefits of Go with the rich ecosystem of AI libraries available in Python.

3. Rudimentary Caching:
   - Due to latency issues with the arXiv API, we implemented a basic caching mechanism. This involves saving a file containing a large number of papers, which we can quickly access instead of making repeated API calls to arXiv. This approach significantly reduces response times and minimizes the load on the arXiv servers.

## System Flow

1. The user interacts with the Angular frontend, which offers two modes: arXiv and Custom.

2. In arXiv mode:
   - The frontend requests random arXiv papers from the Go backend.
   - The Go backend fetches a selection of random papers from its cache or directly from arXiv.
   - The frontend displays the paper titles to the user.
   - When the user selects a paper, the frontend sends an explanation request to the Go backend.

3. In Custom mode:
   - The user inputs their own text for explanation.
   - The frontend sends this text to the Go backend for explanation.

4. For both modes:
   - The Go backend receives the explanation request and calls the Python microservice.
   - The Python microservice interacts with the OpenAI API to generate layered explanations for each concept.
   - The Go backend aggregates the responses and sends the structured explanation back to the frontend.

5. The frontend presents the multi-layered explanation to the user in an interactive format, including:
   - A main takeaway of the entire text.
   - Topics broken down into concepts.
   - Each concept explained in three layers (What, Why, How) of increasing detail.

6. For arXiv papers, the frontend also displays a link to the original PDF, allowing users to access the full paper directly.

This architecture enables PleX to provide fast, scalable, and comprehensive explanations of complex subjects, leveraging the strengths of each technology in our stack while offering flexibility in content sources.

## Features

1. Mode Toggle: Users can switch between 'arXiv' and 'Custom' modes. In arXiv mode, the system fetches a selection of random papers from arXiv to showcase the app's capabilities. In Custom mode, users can input their own text for explanation.

2. PDF Link: For arXiv papers, a link to the original PDF is provided, allowing users to access the full paper directly.

3. Multi-layered Explanations: Each concept is explained in three layers of increasing detail, following the "What, Why, How" framework.

4. Interactive UI: The frontend presents the multi-layered explanation in an interactive format, allowing users to explore concepts at their preferred depth.

## Note on Naming

While the internal project name is Metis, the frontend application is branded as PleX (short for "Please Explain"). This branding decision was made to create a more engaging user experience and to align with the project's goal of making complex ideas simple to understand.

## Future Improvements

1. Optimized arXiv Integration:
   - Implement more efficient caching mechanisms for arXiv papers to reduce API calls and improve response times.
   - Optimize fetching algorithm to reduce latency.

2. Enhanced Caching System:
   - Implement a server-side caching system for explanations. When a user requests an explanation for a specific paper, cache the result for a set period.
   - If another user requests the same paper within the cache period, serve the cached explanation instead of generating a new one, significantly reducing processing time and API usage.

3. Shareable Links:
   - Generate unique, shareable URLs for each explanation.
   - Allow users to easily share interesting explanations with others via a direct link.
   - Implement a system to retrieve cached explanations using these shareable links.

4. Additional Educational Tools:
   - Integrate interactive quizzes and assessments to help users test their understanding of explained concepts.
   - Implement a spaced repetition system to help users retain information over time.
   - Add a feature for users to create and share their own study guides based on the explanations.
   - Incorporate multimedia elements such as diagrams, videos, and interactive simulations to cater to different learning styles.
   - Develop a recommendation system that suggests related topics or papers based on user interests and learning history.

These improvements aim to enhance PleX's performance, user experience, and overall functionality, making it an even more powerful tool for understanding complex ideas and facilitating effective learning across various subjects.