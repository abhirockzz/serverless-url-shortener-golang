package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"url-shortener/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

const queryParameterName = "url"

type Response struct {
	ShortCode string `json:"short_code"`
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	url := req.Body
	log.Println("original url", url)

	shortCode, err := db.SaveURL(url)
	if err != nil {
		log.Println("failed to generate short code for", url)
		return events.APIGatewayV2HTTPResponse{}, err
	}

	response := Response{ShortCode: shortCode}
	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to marshal response for", url)
		return events.APIGatewayV2HTTPResponse{}, err
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusCreated, Body: string(respBytes)}, nil
}
