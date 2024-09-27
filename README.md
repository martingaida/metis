<<<<<<< HEAD
# metis
=======
# Metis: Multi-Layered Concept Explanation System

## Purpose

Metis is an innovative application designed to explain complex subjects by breaking them down into digestible components. The system works by:

1. Dividing the main subject into major topics
2. Breaking down each topic into key concepts
3. Explaining each concept in three layers of increasing detail:
   - Layer 1: High-level overview
   - Layer 2: Detailed explanation
   - Layer 3: Technical deep-dive

Each layer of explanation follows the "What, Why, How" framework to ensure comprehensive understanding:
- What: Definition and basic description of the concept
- Why: Importance and relevance of the concept
- How: Practical application or functioning of the concept

This structured approach allows users to grasp complex information at their preferred depth of understanding, making it suitable for learners at all levels.

## Technical Stack

- Frontend: Angular, Typescript
- Backend: Go
- Microservices: Python
- Deployment: AWS Lambda, S3

### Architecture Decisions

1. Direct Lambda Invocation:
   We've chosen to make direct calls to AWS Lambda functions instead of using API Gateway. This decision was made to reduce latency and costs associated with API Gateway, as our application doesn't require the additional features provided by API Gateway such as request throttling or API key management.

2. Microservices Architecture:
   The backend is implemented in Go for its performance benefits and ease of deployment as a Lambda function. However, the OpenAI client libraries are primarily available for Node.js and Python. To leverage these libraries effectively, we've implemented the LLM (Language Model) interaction as a separate microservice in Python. This architecture allows us to combine the performance benefits of Go with the rich ecosystem of AI libraries available in Python.

## System Flow

1. User inputs a complex subject through the Angular frontend.
2. The Go backend receives the request and orchestrates the explanation process.
3. The Python microservice is called to interact with the OpenAI API, generating layered explanations for each concept.
4. The backend aggregates the responses and sends the structured explanation back to the frontend.
5. The frontend presents the multi-layered explanation to the user in an interactive format.

This architecture enables Metis to provide fast, scalable, and comprehensive explanations of complex subjects, leveraging the strengths of each technology in our stack.
>>>>>>> ddb9400423df2da42dff9107ca11e16756fdb0b7
