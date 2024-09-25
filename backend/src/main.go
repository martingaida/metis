package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: Method=%s, Path=%s, Body=%s", req.HTTPMethod, req.Path, req.Body)

	if req.HTTPMethod == "OPTIONS" {
		return handleOptions(req)
	}

	var request Request
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("Parsed request: %+v", request)

	if request.Text == "" {
		log.Println("Text field is empty")

		return events.APIGatewayProxyResponse{
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

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "LLM_MICROSERVICE_URL not set",
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("Calling LLM microservice at %s", llmServiceURL)
	llmResponse, err := callLLMMicroservice(llmServiceURL, request.Text)
	if err != nil {
		log.Printf("Error calling LLM microservice: %v", err)

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error calling LLM microservice: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("LLM microservice response: %s", llmResponse)

	// Parse LLM response
	var response Response
	err = json.Unmarshal([]byte(llmResponse), &response)
	if err != nil {
		log.Printf("Error parsing LLM response: %v", err)

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error parsing LLM response: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	// Return response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)

		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error marshaling response: " + err.Error(),
			Headers:    getCORSHeaders(),
		}, nil
	}

	log.Printf("Returning response: %s", string(jsonResponse))

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    getCORSHeaders(),
		Body:       string(jsonResponse),
	}, nil
}

func getCORSHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "POST, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
		"Access-Control-Allow-Credentials": "true",
	}
}

func callLLMMicroservice(url, text string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func handleOptions(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "POST, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
