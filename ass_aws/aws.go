package ass_aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewConfig(ctx context.Context, region string) (aws.Config, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)
}
