package vaws

import (
	"context"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awscreds "github.com/aws/aws-sdk-go-v2/credentials"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
)

type AWSConfig struct {
	SecretAccessKey, AccessKeyId, Region string
	KMSKey                               string
}

// getAwsCLientConfig returns the AWS configuration for the given AWSConfig.
// It validates the required fields (Region, AccessKeyId, SecretAccessKey) and returns the AWS configuration.
func (a *AWSConfig) getAwsCLientConfig(ctx context.Context) (aws.Config, error) {
	if a.Region == "" || a.AccessKeyId == "" || a.SecretAccessKey == "" {
		return aws.Config{}, dmerrors.DMError(ErrInvalidAwsConfig, nil)
	}
	return awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion(a.Region),
		awsConfig.WithCredentialsProvider(
			awscreds.NewStaticCredentialsProvider(a.AccessKeyId, a.SecretAccessKey, ""),
		),
	)
}
