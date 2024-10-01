package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	lambdaRuntime "github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Text   string `json:"text"`
	Action string `json:"action"`
	Level  string `json:"level"`
}

type Layer struct {
	What string `json:"what"`
	Why  string `json:"why"`
	How  string `json:"how"`
}

type Concept struct {
	Concept     string `json:"concept"`
	Layer       Layer  `json:"layer"`
	ImagePrompt string `json:"image_prompt"`
	ImageURL    string `json:"image_url"`
}

type Topic struct {
	Topic    string    `json:"topic"`
	Concepts []Concept `json:"concepts"`
}

type Response struct {
	Explanations []Topic `json:"explanations"`
	MainTakeaway string  `json:"main_takeaway"`
}

type LLMResponse struct {
	Topics       []Topic `json:"topics"`
	MainTakeaway string  `json:"main_takeaway"`
}

type ArXivPaper struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Abstract    string `json:"abstract"`
	Category    string `json:"category"`
	Authors     string `json:"authors"`
	Published   string `json:"published"`
	AbstractURL string `json:"abstract_url"`
	PDFURL      string `json:"pdf_url"`
}

func handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Printf("Received request: Method=%s, Path=%s, Body=%s", request.RequestContext.HTTP.Method, request.RequestContext.HTTP.Path, request.Body)

	if request.RequestContext.HTTP.Method == "OPTIONS" {
		return handleOptions()
	}

	// Parse request body
	var requestBody map[string]string
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
			Headers:    getCORSHeaders(),
		}, nil
	}

	// Get the action from the body
	action, exists := requestBody["action"]
	if !exists || action == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Action field is required",
			Headers:    getCORSHeaders(),
		}, nil
	}

	// Handle the action
	switch action {
	case "explain":
		return handleExplain(request)
	case "arxiv":
		return handleGetArXivPapers(request)
	default:
		return events.LambdaFunctionURLResponse{
			StatusCode: 404,
			Body:       "Action not recognized",
			Headers:    getCORSHeaders(),
		}, nil
	}
}

func handleExplain(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	var requestStruct Request
	err := json.Unmarshal([]byte(request.Body), &requestStruct)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Invalid request body: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	if requestStruct.Text == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Text field is required",
			Headers:    getCORSHeaders(),
		}, nil
	}

	if requestStruct.Level == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Level field is required",
			Headers:    getCORSHeaders(),
		}, nil
	}

	// Call LLM microservice
	llmResponse, err := callLLMMicroservice(requestStruct.Text, requestStruct.Level, false)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error calling LLM microservice: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	// Prepare the response
	response := Response{
		Explanations: llmResponse.Topics,
		MainTakeaway: llmResponse.MainTakeaway,
	}

	responseJSON, err := json.Marshal(response)
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

func handleGetArXivPapers(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// Call the ArXiv microservice
	arxivResponse, err := callArXivMicroservice()
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error calling ArXiv microservice: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, err
	}

	responseJSON, err := json.Marshal(arxivResponse)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error marshaling ArXiv response: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, err
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

func callLLMMicroservice(text string, level string, generateImages bool) (LLMResponse, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"text":            text,
		"level":           level,
		"generate_images": generateImages,
	})
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

func callArXivMicroservice() ([]ArXivPaper, error) {
	arxivURL := os.Getenv("ARXIV_MICROSERVICE_URL")
	if arxivURL == "" {
		return nil, fmt.Errorf("ARXIV_MICROSERVICE_URL is not set")
	}

	log.Printf("Calling ArXiv microservice at: %s", arxivURL)

	resp, err := http.Get(arxivURL)
	if err != nil {
		return nil, fmt.Errorf("error calling ArXiv microservice: %v", err)
	}
	defer resp.Body.Close()

	var papers []ArXivPaper
	err = json.NewDecoder(resp.Body).Decode(&papers)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return papers, nil
}

func main() {
	lambdaRuntime.Start(handleRequest)
}
