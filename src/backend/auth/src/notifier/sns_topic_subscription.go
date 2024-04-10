package notifier

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
)

func AddSNSSubscription(snsClient *sns.SNS, sns_arn, email string) {

	filterPolicy := map[string]interface{}{
		"target_email": email,
	}

	filterPolicyJson, jerr := json.Marshal(filterPolicy)
	if jerr != nil {
		fmt.Println("Error marshaling filter policy:", jerr)
		return
	}

	snsInput := &sns.SubscribeInput{
		TopicArn: aws.String(sns_arn),
		Protocol: aws.String("email"),
		Endpoint: aws.String(email),
		Attributes: map[string]*string{
			"FilterPolicy": aws.String(string(filterPolicyJson)),
		},
	}
	_, err := snsClient.Subscribe(snsInput)
	log.Printf("Successfully created SNS subscription for %s", email)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sns.ErrCodeSubscriptionLimitExceededException:
				log.Printf("The subscription limit has been reached for the SNS topic.")
				os.Exit(0)
			case sns.ErrCodeInvalidParameterException:
				log.Printf("Invalid parameter to create SNS subscription")
				os.Exit(0)
			case sns.ErrCodeNotFoundException:
				log.Printf("The requested resource(SNS_TOPIC) does not exist")
				os.Exit(0)
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}
