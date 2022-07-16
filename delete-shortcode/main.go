package main

import (
	"context"
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

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	shortCode := req.PathParameters[pathParameterName]

	log.Println("delete request for short-code", shortCode)

	err := db.Delete(shortCode)
	if err != nil {
		if errors.Is(err, db.ErrUrlNotFound) {
			return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound, Body: fmt.Sprintf("short code %s not found\n", shortCode)}, nil
		}
		return events.APIGatewayV2HTTPResponse{}, err
	}
	log.Println("successfully deleted", shortCode)

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNoContent}, nil
}
