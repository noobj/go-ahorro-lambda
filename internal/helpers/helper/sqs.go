package helper

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendSqsMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	qURL := os.Getenv("SQS_URL")
	input.QueueUrl = &qURL

	result, err := svc.SendMessage(input)

	if err != nil {
		return nil, err
	}

	return result, nil
}
