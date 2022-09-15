package helper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/noobj/go-serverless-services/internal/config"
)

func SendSqsMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	env := config.GetInstance()
	svc := sqs.New(sess)
	qURL := env.SqsUrl
	input.QueueUrl = &qURL

	result, err := svc.SendMessage(input)

	if err != nil {
		return nil, err
	}

	return result, nil
}
