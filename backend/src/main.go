package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Text string `json:"text"`
}

type Layer struct {
	LayerName   string `json:"layer_name"`
	Explanation string `json:"explanation"`
}

type Concept struct {
	Concept string  `json:"concept"`
	Layers  []Layer `json:"layers"`
}

type Topic struct {
	Topic    string    `json:"topic"`
	Concepts []Concept `json:"concepts"`
}

type Response struct {
	Explanations []Topic `json:"explanations"`
}

func handleRequest(req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Printf("Received request: Method=%s, Body=%s", req.RequestContext.HTTP.Method, req.Body)

	if req.RequestContext.HTTP.Method == "OPTIONS" {
		return handleOptions()
	}

	var request Request
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Invalid request body: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("Parsed request: %+v", request)

	if request.Text == "" {
		log.Println("Text field is empty")
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Text field is required",
			Headers:    getCORSHeaders(),
		}, nil
	}

	// Call LLM microservice
	llmServiceURL := os.Getenv("LLM_MICROSERVICE_URL")
	log.Printf("LLM_MICROSERVICE_URL: %s", llmServiceURL)
	if llmServiceURL == "" {
		log.Println("LLM_MICROSERVICE_URL not set")
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "LLM_MICROSERVICE_URL not set",
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("Calling LLM microservice at %s", llmServiceURL)
	llmResponse, err := callLLMMicroservice(llmServiceURL, request.Text)
	if err != nil {
		log.Printf("Error calling LLM microservice: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error calling LLM microservice: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("LLM microservice response: %s", llmResponse)

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       llmResponse,
		Headers:    getCORSHeaders(),
	}, nil
}

func handleOptions() (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    getCORSHeaders(),
	}, nil
}

func getCORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}
}

func callLLMMicroservice(url, text string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Message   string `json:"message"`
		RequestID string `json:"requestId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Here, instead of polling, we're just returning a message that processing has started
	// In a real-world scenario, you might implement a way to check the status or retrieve the result later
	return fmt.Sprintf("Processing started. Request ID: %s", response.RequestID), nil
}

func main() {
	lambda.Start(handleRequest)
}
