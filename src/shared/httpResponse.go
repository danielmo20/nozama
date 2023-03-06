package nozama

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func HttpResponse(statusCode int, obj interface{}) (events.APIGatewayProxyResponse, error) {

	if obj == nil {
		obj = HttpMessageResponse{http.StatusText(statusCode)}
	}

	objBytes, err := json.Marshal(obj)

	if err != nil {
		log.Printf("HttpResponse: Cannont json.Marshal this %v", objBytes)
	}

	message := string(objBytes)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       message,
	}, nil

}
