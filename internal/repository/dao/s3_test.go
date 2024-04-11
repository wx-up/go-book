package dao

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Test(t *testing.T) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "oss",
			URL:           "https://oss-cn-hongkong.aliyuncs.com",
			SigningRegion: "cn-hongkong",
		}, nil
	})
	accessKeyID := os.Getenv("OSS_AccessKeyId")
	accessKeySecret := os.Getenv("OSS_AccessKeySecret")
	cred := credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(cred),
		config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      nil,
		Key:         nil,
		Body:        nil,
		ContentType: nil,
	})
	assert.NoError(t, err)
}
