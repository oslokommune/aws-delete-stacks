package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/oslokommune/aws-delete-stacks/cmd"
	"github.com/oslokommune/aws-delete-stacks/core/application"
	"github.com/oslokommune/aws-delete-stacks/core/repository"
	"os"
)

func main() {
	awsSession := newAwsSession()
	awsCloudFormationRepository := repository.NewAWSCloudFormationRepository(awsSession)
	cloudFormation := application.NewCloudFormation(awsCloudFormationRepository)

	deleteCmd := cmd.BuildDeleteCommand(cloudFormation)

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
