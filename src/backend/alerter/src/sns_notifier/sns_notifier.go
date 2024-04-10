package sns_notifier

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
)

type MessageAttributes struct {
	Message string
	Body    string
}

func SnsMessageFormat(username, vendor string) MessageAttributes {
	/*
		Formats the email body and message to be sent
		out for expiring subscriptions
		Params: username string
				vendor string
		Return: MessageAttributes
	*/

	subject := fmt.Sprintf("Important Notice: Your %s Subscription Renewal", vendor)

	message := fmt.Sprintf(`Dear %s,

Your subscription to %s is set to expire tomorrow. We want to make sure you have the opportunity to continue your service or make changes if needed.

If you have any questions or require assistance, please contact our customer support team.
	
Thank you for being a valued customer.
	
Sincerely,
SUBHUB
	`, username, vendor)

	return MessageAttributes{
		Message: subject,
		Body:    message,
	}
}

func PublishMessage(snsCli *sns.SNS, snsArn, subject, body, email string) {
	/*
		Publishes an SNS message to the specified email
		address
		Params: snsCli *sns.SNS,
				snsArn string
				subject string
				body string
				email string
		Return: None
	*/
	input := &sns.PublishInput{
		TopicArn: aws.String(snsArn),
		Subject:  aws.String(subject),
		Message:  aws.String(body),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"target_email": {
				DataType:    aws.String("String"),
				StringValue: aws.String(email),
			},
		},
	}
	fmt.Printf("Sending alert to email %s", email)
	_, err := snsCli.Publish(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sns.ErrCodeInvalidParameterException, sns.ErrCodeValidationException:
				log.Printf("Invalid parameter for subscription for email %s", email)
				os.Exit(1)
			case sns.ErrCodeNotFoundException:
				log.Printf("Unable to find SNS topic for pushing subscription, please check arn %s", snsArn)
				os.Exit(1)
			case sns.ErrCodeEndpointDisabledException:
				log.Printf("Invalid Endpoint for subscribing email %s", email)
				os.Exit(1)
			}
		} else {
			log.Printf("Got error calling GetItem: %s", err)
			os.Exit(1)
		}
	}
	if err != nil {
		fmt.Printf("Unable to Publish message to %s", email)
		os.Exit(1)
	}
}
