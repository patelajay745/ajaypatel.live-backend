package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	sesClient *ses.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig((context.TODO()))
	if err != nil {
		fmt.Printf("unable to load SDK config,%v", err)
	}
	sesClient = ses.NewFromConfig(cfg)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println(request.Body)

	values, err := url.ParseQuery(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error parsing form data",
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "https://ajaypatel.live",
				"Access-Control-Allow-Methods": "POST, GET",
			},
		}, err
	}

	if values.Get("firstName") == "" ||
		values.Get("lastName") == "" ||
		values.Get("email") == "" ||
		values.Get("subject") == "" ||
		values.Get("message") == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Please provide all required fields",
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "https://ajaypatel.live",
				"Access-Control-Allow-Methods": "POST, GET",
			},
		}, nil
	}

	message := values.Get("message") + " From : " + values.Get("email")

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				"patel.ajay745@gmail.com",
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(message),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(values.Get("subject")),
			},
		},
		Source: aws.String("patel.ajay745@gmail.com"),
	}

	result, err := sesClient.SendEmail(ctx, input)

	if err != nil {
		fmt.Println("Error sending email:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error sending email",
			Headers: map[string]string{
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Origin":  "https://ajaypatel.live",
				"Access-Control-Allow-Methods": "POST, GET",
			},
		}, err
	}
	fmt.Println("Email sent! Message ID:", result.MessageId)

	return events.APIGatewayProxyResponse{
		Body:       "Email sent",
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "https://ajaypatel.live",
			"Access-Control-Allow-Methods": "POST, GET",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
