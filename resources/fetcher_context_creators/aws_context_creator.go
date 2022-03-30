package fetcher_context_creators

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"log"
)

const (
	Name = "aws_context_creator"
	ContextName = "aws_context"
)

type awsConfigContextCreator struct{}

func (cr *awsConfigContextCreator) GetContext() (interface{}, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	return cfg, err
}

func (cr *awsConfigContextCreator) GetName() string {
	return Name
}

func (cr *awsConfigContextCreator) GetContextName() string {
	return ContextName
}
