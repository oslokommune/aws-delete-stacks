package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/oslokommune/aws-delete-stacks/aws_cloudformation"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"os"
)

func main() {
	awsSession := newAwsSession()
	awsCloudFormation := aws_cloudformation.NewAWSCloudFormation(awsSession)
	deleteCmd := cmd.BuildDeleteCommand(awsCloudFormation)

	if err := deleteCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newAwsSession() *session.Session {
	t := true
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: &t,
		},
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}
