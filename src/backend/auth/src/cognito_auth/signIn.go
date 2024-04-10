package cognito_auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log/slog"
	"os"
)

const (
	USERNAME = "USERNAME"
	PASSWORD = "PASSWORD"
)

func SignIn(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider, client_id, username, password string) ReturnResults {

	var res ReturnResults // ????
	jsonhandler := slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(jsonhandler)

	input := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(client_id),
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			USERNAME: aws.String(username),
			PASSWORD: aws.String(password),
		},
		ClientMetadata: map[string]*string{
			USERNAME: aws.String(username),
			PASSWORD: aws.String(password),
		},
	}
	_, err := cognitoClient.InitiateAuth(input)
	logger.Info("Successfully Authenticated user " + username)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodePasswordResetRequiredException:
				logger.Error(username + "needs to reset password")
				res.Message = fmt.Sprintf(username + "needs to reset password")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				logger.Error("User" + username + "does not exist")
				res.Message = fmt.Sprintf("User" + username + "does not exist")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeUserNotConfirmedException:
				logger.Error(username + "needs to confirm account")
				res.Message = fmt.Sprintf(username + "needs to confirm account")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeNotAuthorizedException:
				logger.Error(username + "is not authorized")
				res.Message = fmt.Sprintf(username + "is not authorized to perform this action")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeLimitExceededException:
				logger.Error(username + "has exceeded the no of tries for logging in")
				res.Message = fmt.Sprintf("has exceeded the no of tries for logging in")
				res.Status = 400
				return res
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	res.Message = fmt.Sprintf("Succesfully signed in user %s", username)
	res.Status = 200

	return res
}
