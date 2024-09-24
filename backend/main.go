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

type Response struct {
	Explanation string `json:"explanation"`
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

	response := Response{Explanation: llmResponse}
	jsonResponse, _ := json.Marshal(response)

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

func callLLMMicroservice(url, text string) (string, error) {
	requestBody, _ := json.Marshal(map[string]string{"text": text})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return result["explanation"].(string), nil
}

func main() {
	lambda.Start(handleRequest)
}