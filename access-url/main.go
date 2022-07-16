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

const locationHeader = "Location"
const pathParameterName = "shortcode"

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	shortCode := req.PathParameters[pathParameterName]

	log.Println("redirect request for short-code", shortCode)

	longurl, err := db.GetLongURL(shortCode)

	if err != nil {
		if errors.Is(err, db.ErrUrlNotFound) {
			return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound, Body: fmt.Sprintf("short code %s not found\n", shortCode)}, nil
		} else if errors.Is(err, db.ErrUrlNotActive) {
			return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusBadRequest, Body: fmt.Sprintf("short code %s not active\n", shortCode)}, nil
		}
		return events.APIGatewayV2HTTPResponse{}, err
	}

	log.Println("redirecting to ", longurl)

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusFound, Headers: map[string]string{locationHeader: longurl}}, nil
}
