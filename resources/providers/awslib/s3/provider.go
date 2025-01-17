// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package s3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	s3Client "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/elastic/cloudbeat/resources/fetching"
	"github.com/elastic/cloudbeat/resources/providers/awslib"
	"github.com/elastic/elastic-agent-libs/logp"
)

func NewProvider(cfg aws.Config, log *logp.Logger, factory awslib.CrossRegionFactory[Client]) *Provider {
	f := func(cfg aws.Config) Client {
		return s3Client.NewFromConfig(cfg)
	}
	m := factory.NewMultiRegionClients(ec2.NewFromConfig(cfg), cfg, f, log)

	return &Provider{
		log:     log,
		clients: m.GetMultiRegionsClientMap(),
	}
}

func (p Provider) DescribeBuckets(ctx context.Context) ([]awslib.AwsResource, error) {
	clientBuckets, err := p.clients[awslib.DefaultRegion].ListBuckets(ctx, &s3Client.ListBucketsInput{})
	if err != nil {
		p.log.Errorf("Could not list s3 buckets: %v", err)
		return nil, err
	}

	var result []awslib.AwsResource
	bucketsRegionsMapping := p.getBucketsRegionMapping(ctx, clientBuckets.Buckets)
	for region, buckets := range bucketsRegionsMapping {
		for _, bucket := range buckets {
			// Getting the bucket encryption, policy and versioning is not critical for the rest of the flow, so we
			//	should keep describing the bucket even if getting these objects fails.
			sseAlgorithm, encryptionErr := p.getBucketEncryptionAlgorithm(ctx, bucket.Name, region)
			if encryptionErr != nil {
				p.log.Errorf("Could not get encryption for bucket %s. Error: %v", *bucket.Name, encryptionErr)
			}

			bucketPolicy, policyErr := p.GetBucketPolicy(ctx, bucket.Name, region)

			if policyErr != nil {
				p.log.Errorf("Could not get bucket policy for bucket %s. Error: %v", *bucket.Name, policyErr)
			}

			bucketVersioning, versioningErr := p.getBucketVersioning(ctx, bucket.Name, region)
			if versioningErr != nil {
				p.log.Errorf("Could not get bucket versioning for bucket %s. Err: %v", *bucket.Name, versioningErr)
			}

			result = append(result, BucketDescription{*bucket.Name, sseAlgorithm, bucketPolicy, bucketVersioning})
		}
	}

	return result, nil
}

func (p Provider) GetBucketACL(ctx context.Context, bucketName *string, region string) (*s3Client.GetBucketAclOutput, error) {
	client, err := awslib.GetClient(&region, p.clients)
	if err != nil {
		return nil, err
	}

	acl, err := client.GetBucketAcl(ctx, &s3Client.GetBucketAclInput{Bucket: bucketName})
	if err != nil {
		p.log.Debugf("Error getting bucket ACL: %s", err)
		return nil, err
	}

	return acl, nil
}

func (p Provider) GetBucketPolicy(ctx context.Context, bucketName *string, region string) (BucketPolicy, error) {
	client, err := awslib.GetClient(&region, p.clients)
	if err != nil {
		return nil, err
	}

	rawPolicy, err := client.GetBucketPolicy(ctx, &s3Client.GetBucketPolicyInput{Bucket: bucketName})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchBucketPolicy" {
				p.log.Debugf("Bucket policy for bucket %s does not exist", *bucketName)
				return nil, nil
			}
		}

		return nil, err
	}

	var bucketPolicy BucketPolicy
	jsonErr := json.Unmarshal([]byte(*rawPolicy.Policy), &bucketPolicy)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return bucketPolicy, nil
}

func (p Provider) GetBucketLogging(ctx context.Context, bucketName *string, region string) (Logging, error) {
	client, err := awslib.GetClient(&region, p.clients)
	if err != nil {
		return Logging{}, err
	}

	logging, err := client.GetBucketLogging(ctx, &s3Client.GetBucketLoggingInput{Bucket: bucketName})
	if err != nil {
		p.log.Debugf("Error getting bucket logging: %s", err)
		return Logging{}, err
	}

	bucketLogging := Logging{}
	if logging.LoggingEnabled != nil {
		bucketLogging.Enabled = true
		bucketLogging.TargetBucket = *logging.LoggingEnabled.TargetBucket
	}

	return bucketLogging, nil
}

func (p Provider) getBucketsRegionMapping(ctx context.Context, buckets []types.Bucket) map[string][]types.Bucket {
	bucketsRegionMap := make(map[string][]types.Bucket, 0)
	for _, clientBucket := range buckets {
		region, regionErr := p.getBucketRegion(ctx, clientBucket.Name)
		// If we could not get the region for a bucket, additional API calls for resources will probably fail, we should
		//	not describe this bucket.
		if regionErr != nil {
			p.log.Errorf("Could not get bucket location for bucket %s. Not describing this bucket. Error: %v", *clientBucket.Name, regionErr)
			continue
		}

		bucketsRegionMap[region] = append(bucketsRegionMap[region], clientBucket)
	}

	return bucketsRegionMap
}

func (p Provider) getBucketEncryptionAlgorithm(ctx context.Context, bucketName *string, region string) (string, error) {
	client, err := awslib.GetClient(&region, p.clients)
	if err != nil {
		return "", err
	}

	encryption, err := client.GetBucketEncryption(ctx, &s3Client.GetBucketEncryptionInput{Bucket: bucketName})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "ServerSideEncryptionConfigurationNotFoundError" {
				p.log.Debugf("Bucket encryption for bucket %s does not exist", *bucketName)
				return "", nil
			}
		}

		return "", err
	}

	if len(encryption.ServerSideEncryptionConfiguration.Rules) <= 0 {
		return "", nil
	}

	return string(encryption.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm), nil
}

func (p Provider) getBucketRegion(ctx context.Context, bucketName *string) (string, error) {
	location, err := p.clients[awslib.DefaultRegion].GetBucketLocation(ctx, &s3Client.GetBucketLocationInput{Bucket: bucketName})
	if err != nil {
		return "", err
	}

	region := string(location.LocationConstraint)
	// Region us-east-1 have a LocationConstraint of null.
	if region == "" {
		region = "us-east-1"
	}

	return region, nil
}

func (p Provider) getBucketVersioning(ctx context.Context, bucketName *string, region string) (BucketVersioning, error) {
	bucketVersioning := BucketVersioning{false, false}

	client, err := awslib.GetClient(&region, p.clients)
	if err != nil {
		return bucketVersioning, err
	}

	bucketVersioningResponse, err := client.GetBucketVersioning(ctx, &s3Client.GetBucketVersioningInput{Bucket: bucketName})
	if err != nil {
		return bucketVersioning, err
	}

	if bucketVersioningResponse.Status == types.BucketVersioningStatusEnabled {
		bucketVersioning.Enabled = true
	}

	if bucketVersioningResponse.MFADelete == types.MFADeleteStatusEnabled {
		bucketVersioning.MfaDelete = true
	}

	return bucketVersioning, nil
}

func (b BucketDescription) GetResourceArn() string {
	return fmt.Sprintf("arn:aws:s3:::%s", b.Name)
}

func (b BucketDescription) GetResourceName() string {
	return b.Name
}

func (b BucketDescription) GetResourceType() string {
	return fetching.S3Type
}
