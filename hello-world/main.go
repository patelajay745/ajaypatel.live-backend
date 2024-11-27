package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SenderInformation struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

type Response struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

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
		}, err
	}

	senderInfo := SenderInformation{
		FirstName: values.Get("firstName"),
		LastName:  values.Get("lastName"),
		Email:     values.Get("email"),
		Subject:   values.Get("subject"),
		Message:   values.Get("message"),
	}

	const (
		Sender    = values.Get("email")
		Recipient = "patel.ajay745@gmail.com"
		Subject   = values.Get("subject")
		Body      = values.Get("message")
		CharSet   = "UTF-8"
	)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	svc := ses.New(sess)

	// var senderInfo SenderInformation
	// if err := json.Unmarshal([]byte(request.Body), &senderInfo); err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 400,
	// 		Body:       "Error parsing sender information",
	// 	}, err
	// }

	// fmt.Println("This is Form-Data", request.Body)

	// response ={
	// 	"firstname": senderInfo.FirstName,
	// }

	return events.APIGatewayProxyResponse{
		Body:       senderInfo.FirstName,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
