package fetchers

import (
	"encoding/gob"
	"github.com/elastic/cloudbeat/resources/ctxProvider"

	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/cloudbeat/resources/fetching"
	"github.com/elastic/cloudbeat/resources/manager"
)

const (
	IAMType = "aws-iam"
)

func init() {
	manager.Factories.ListFetcherFactory(IAMType, &IAMFactory{})
	gob.Register(IAMResource{})
}

type IAMFactory struct {
}

func (f *IAMFactory) Create(c *common.Config) (fetching.Fetcher, error) {
	cfg := IAMFetcherConfig{}
	err := c.Unpack(&cfg)
	if err != nil {
		return nil, err
	}

	return f.CreateFrom(cfg)
}

func (f *IAMFactory) CreateFrom(cfg IAMFetcherConfig) (fetching.Fetcher, error) {
	awsCredProvider := ctxProvider.AWSCredProvider{}
	awsCfg := awsCredProvider.GetAwsCredentials()
	iam := NewIAMProvider(awsCfg.Config)

	return &IAMFetcher{
		cfg:         cfg,
		iamProvider: iam,
	}, nil

}
