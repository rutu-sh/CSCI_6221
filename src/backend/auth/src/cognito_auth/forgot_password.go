package cognito_auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log/slog"
	"os"
)

func ForgotPassword(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider, client_id, username string) ReturnResults {

	jsonhandler := slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(jsonhandler)
	var res ReturnResults // ????

	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(client_id),
		Username: aws.String(username),
	}
	_, err := cognitoClient.ForgotPassword(input)
	logger.Info("Successfully Authenticated user " + username)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				logger.Error("User" + username + "does not exist")
				res.Message = fmt.Sprintf("User" + username + "does not exist")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeCodeMismatchException:
				logger.Error("Invalid Verification code")
				res.Message = fmt.Sprintf("Invalid Verification code")
				res.Status = 400
				return res
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	res.Message = fmt.Sprintf("Succesfully resent verification code to %s", username)
	res.Status = 200

	return res
}
