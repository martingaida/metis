<p align="center">
  <img src="frontend/src/assets/images/plex-logo.svg" alt="PleX Logo" width="220">
</p>

# PleX: Complex Ideas Made Simple

PleX (short for "Please Explain") is an innovative application designed to explain complex subjects by adapting explanations to different reading levels. The system works by:

1. Analyzing the input text or topic
2. Breaking down the subject into key concepts
3. Explaining each concept at the user's chosen reading level, based on the Flesch-Kincaid Scale:
   - K3: Ages 5-8, equivalent to kindergarten to 3rd grade
   - K6: Ages 8-11, equivalent to 4th to 6th grade
   - K9: Ages 11-14, equivalent to 7th to 9th grade
   - K12: Ages 14-17, equivalent to 10th to 12th grade
   - College: Ages 17-20, equivalent to undergraduate level
   - Graduate: Ages 20+, equivalent to graduate or professional level
4. Each explanation follows the "What, Why, How" framework to ensure comprehensive understanding:
   - What: Definition and basic description of the concept
   - Why: Importance and relevance of the concept
   - How: Practical application or functioning of the concept

This adaptive approach allows users to grasp complex information at their preferred level of understanding, making it suitable for learners of all ages and backgrounds.

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

5. Image Generation Feature

   - Our application now includes an exciting image generation feature that enhances the explanation process by providing visual representations of key concepts. This feature utilizes state-of-the-art AI technology to create relevant images for each concept explained.

   - **Model**: We use the DALL-E 3 model, which is known for its ability to generate high-quality, contextually relevant images based on text prompts.
   
   - **Integration**: The image generation is seamlessly integrated into our explanation pipeline. After generating the textual explanation for each concept, we automatically create a prompt for image generation.
   
   - **Image Specifications**: 
      1. Size: 1024x1024 pixels
      2. Quality: Standard
      3. Number of images per concept: 1

### How it works:

1. When an explanation is generated, each concept within the explanation is processed for image creation.
2. An image prompt is automatically generated based on the concept's description.
3. This prompt is sent to the DALL-E 3 model via the OpenAI API.
4. The generated image URL is then included in the response alongside the textual explanation.
5. In the frontend, these images are displayed next to their corresponding concepts, providing a visual aid to the explanation.

### Benefits:

- Enhanced Understanding: Visual representations can significantly improve comprehension of complex concepts.
- Engagement: Images make the learning experience more engaging and memorable.
- Accessibility: Visual aids can help users with different learning styles or those who prefer visual information.

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