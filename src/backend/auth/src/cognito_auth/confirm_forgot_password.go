package cognito_auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log/slog"
	"os"
)

func ConfirmForgetPassword(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider, client_id, username, confCode, newPassword string) ReturnResults {

	jsonhandler := slog.NewJSONHandler(os.Stderr, nil)
	logger := slog.New(jsonhandler)
	var res ReturnResults

	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(client_id),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confCode),
		Password:         aws.String(newPassword),
	}
	_, err := cognitoClient.ConfirmForgotPassword(input)
	logger.Info("Successfully Authenticated user " + username)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				logger.Error("User" + username + "does not exist")
				res.Message = fmt.Sprintf("User has exceeded the OTP entry limit")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeCodeMismatchException:
				logger.Error("User" + username + " provided Invalid Verification code")
				res.Message = fmt.Sprintf("User has exceeded the OTP entry limit")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeTooManyRequestsException:
				logger.Error("User" + username + "has exceeded the retry quota")
				res.Message = fmt.Sprintf("User has exceeded the OTP entry limit")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeLimitExceededException:
				logger.Error("User" + username + "has exceeded the request limit")
				res.Message = fmt.Sprintf("User has exceeded the OTP entry limit")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeTooManyFailedAttemptsException:
				logger.Error("User" + username + "has exceeded the number of attempts")
				res.Message = fmt.Sprintf("User has exceeded the OTP entry limit")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeUserNotConfirmedException:
				logger.Error("User" + username + "has not confirmed email")
				res.Message = fmt.Sprintf("User has not confirmed their email")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeExpiredCodeException:
				logger.Error("User" + username + "provided OTP has expired")
				res.Message = fmt.Sprintf("User has provided an expired OTP")
				res.Status = 400
				return res

			}
		} else {
			fmt.Println(err.Error())
		}
	}
	res.Message = fmt.Sprintf("Succesfully changed password for %s", username)
	res.Status = 200

	return res
}
