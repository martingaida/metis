package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "bytes"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

// LLMRequest struct to hold incoming text
type LLMRequest struct {
    Text string `json:"text"`
}

// LLMResponse struct to hold outgoing response
type LLMResponse struct {
    Result string `json:"result"`
}

// Handler function for AWS Lambda
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var llmRequest LLMRequest
    err := json.Unmarshal([]byte(request.Body), &llmRequest)
    if err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: 400,
            Body:       "Invalid input",
        }, nil
    }

    // Call the LLM microservice to process the text
    llmResponse, err := callMicroserviceLLM(llmRequest.Text)
    if err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: 500,
            Body:       "Error from LLM microservice",
        }, nil
    }

    // Convert the response to JSON
    responseBody, _ := json.Marshal(llmResponse)
    
    // Return the response to the API Gateway
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       string(responseBody),
    }, nil
}

// Function to call the LLM microservice
func callMicroserviceLLM(text string) (LLMResponse, error) {
    llmURL := os.Getenv("LLM_MICROSERVICE_URL")
    if llmURL == "" {
        return LLMResponse{}, fmt.Errorf("LLM_MICROSERVICE_URL environment variable not set")
    }

    requestBody, err := json.Marshal(LLMRequest{Text: text})
    if err != nil {
        return LLMResponse{}, err
    }

    resp, err := http.Post(llmURL, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return LLMResponse{}, err
    }
    defer resp.Body.Close()

    var llmResponse LLMResponse
    err = json.NewDecoder(resp.Body).Decode(&llmResponse)
    if err != nil {
        return LLMResponse{}, err
    }

    return llmResponse, nil
}

func main() {
    // Start the Lambda function handler
    lambda.Start(handler)
}