package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"url-shortener/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

const pathParameterName = "shortcode"

type Payload struct {
	Active bool `json:"active"`
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	var payload Payload
	reqBody := req.Body

	err := json.Unmarshal([]byte(reqBody), &payload)
	if err != nil {
		log.Println("failed to unmarshal status", err)
		return events.APIGatewayV2HTTPResponse{}, err
	}

	shortCode := req.PathParameters[pathParameterName]
	log.Println("status update request for short-code", shortCode, payload.Active)

	err = db.Update(shortCode, payload.Active)
	if err != nil {
		if errors.Is(err, db.ErrUrlNotFound) {
			return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound, Body: fmt.Sprintf("short code %s not found\n", shortCode)}, nil
		}
		return events.APIGatewayV2HTTPResponse{}, err
	}

	log.Println("successfully updated status for", shortCode)

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNoContent}, nil
}
