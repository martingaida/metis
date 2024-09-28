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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type Request struct {
	Text string `json:"text"`
}

type Layer struct {
	What string `json:"what"`
	Why  string `json:"why"`
	How  string `json:"how"`
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
	Explanations struct {
		Topics       []Topic `json:"topics"`
		MainTakeaway string  `json:"main_takeaway"`
	} `json:"explanations"`
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

	switch request.RequestContext.HTTP.Path {
	case "/api/explain":
		return handleExplain(request)
	case "/api/arxiv":
		return handleGetArXivPapers(request)
	default:
		return events.LambdaFunctionURLResponse{
			StatusCode: 404,
			Body:       "Not Found",
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

	llmResponse, err := callLLMMicroservice(requestStruct.Text)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error calling LLM microservice: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

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

func handleGetArXivPapers(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	papers, err := callArXivMicroservice()
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error calling ArXiv microservice: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	responseJSON, err := json.Marshal(papers)
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

func callArXivMicroservice() ([]ArXivPaper, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2")) // Replace with your AWS region
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	payload, err := json.Marshal(map[string]string{})
	if err != nil {
		return nil, err
	}

	result, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: aws.String("ArXivFunction"), // Replace with your actual function name
		Payload:      payload,
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		StatusCode int             `json:"statusCode"`
		Body       json.RawMessage `json:"body"`
	}
	err = json.Unmarshal(result.Payload, &response)
	if err != nil {
		return nil, err
	}

	var papers []ArXivPaper
	err = json.Unmarshal(response.Body, &papers)
	if err != nil {
		return nil, err
	}

	return papers, nil
}

func main() {
	lambdaRuntime.Start(handleRequest)
}
