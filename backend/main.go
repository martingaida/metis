package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request Request
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
	}

	llmServiceURL := os.Getenv("LLM_MICROSERVICE_URL")
	if llmServiceURL == "" {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "LLM_MICROSERVICE_URL not set"}, nil
	}

	// Call LLM microservice
	llmResponse, err := callLLMMicroservice(llmServiceURL, request.Text)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error calling LLM microservice: %v", err)}, nil
	}

	// Preserve the structure by directly passing the llmResponse
	jsonResponse, _ := json.Marshal(llmResponse)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
		Body: string(jsonResponse),
	}, nil
}

func handleOptions(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}

func callLLMMicroservice(url, text string) (Response, error) {
	requestBody, _ := json.Marshal(map[string]string{"text": text})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Response{}, err
	}

	return result, nil
}

func main() {
	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if req.HTTPMethod == "OPTIONS" {
			return handleOptions(req)
		}
		return handleRequest(req)
	})
}
