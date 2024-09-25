package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

type LLMResponse struct {
	Explanations []Topic `json:"explanations"`
}

func handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Printf("Received request: Method=%s, Body=%s", request.RequestContext.HTTP.Method, request.Body)

	if request.RequestContext.HTTP.Method == "OPTIONS" {
		return handleOptions()
	}

	var requestStruct Request
	err := json.Unmarshal([]byte(request.Body), &requestStruct)
	if err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Invalid request body: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("Parsed request: %+v", requestStruct)

	if requestStruct.Text == "" {
		log.Println("Text field is empty")
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Text field is required",
			Headers:    getCORSHeaders(),
		}, nil
	}

	// Call LLM microservice
	llmServiceURL := os.Getenv("LLM_MICROSERVICE_URL")

	if llmServiceURL == "" {
		log.Println("LLM_MICROSERVICE_URL not set")
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "LLM_MICROSERVICE_URL not set",
			Headers:    getCORSHeaders(),
		}, nil
	}

	llmResponse, err := callLLMMicroservice(llmServiceURL + "?text=" + url.QueryEscape(requestStruct.Text))

	if err != nil {
		log.Printf("Error calling LLM microservice: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error calling LLM microservice: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("LLM microservice response: %s", llmResponse)

	responseJSON, err := json.Marshal(llmResponse)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error marshaling response: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
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

func callLLMMicroservice(text string) (LLMResponse, error) {
	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return LLMResponse{}, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(os.Getenv("LLM_MICROSERVICE_URL"), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return LLMResponse{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	var llmResponse LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&llmResponse)
	if err != nil {
		return LLMResponse{}, fmt.Errorf("error decoding response: %v", err)
	}

	return llmResponse, nil
}

func main() {
	lambda.Start(handleRequest)
}
