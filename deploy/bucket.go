package deploy

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CreateBucket creates the bucket from where the website will be hosted
func CreateBucket(config *StaticWebConfig, s3Session *s3.S3) error {

	_, err := s3Session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(config.DomainName),
	})

	if err != nil {
		fmt.Printf("Unable to create bucket %q, %v", config.DomainName, err)
		return err
	}

	s3Session.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(config.DomainName),
	})

	log.Printf("Bucket %s has been created", config.DomainName)

	return nil
}

// PolicyJson gives public read access to a bucket
const PolicyJson = "{\"Version\": \"2008-10-17\",\"Id\": \"PolicyForPublicWebsiteContent\",\"Statement\": [{\"Sid\": \"PublicReadGetObject\",\"Effect\": \"Allow\",\"Principal\": {\"AWS\": \"*\"},\"Action\": \"s3:GetObject\",\"Resource\": \"arn:aws:s3:::[BUCKETNAMEHERE]/*\"}]}"

func SetBucketPermissions(config *StaticWebConfig, s3Session *s3.S3) error {

	_, err := s3Session.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(config.DomainName),
		Policy: aws.String(strings.Replace(PolicyJson, "[BUCKETNAMEHERE]", config.DomainName, 1)),
	})

	if err != nil {
		fmt.Printf("Unable to update policy bucket %s, %v", config.DomainName, err)
		return err
	}

	return nil
}

func CreateBucketWebsite(config *StaticWebConfig, s3Session *s3.S3) error {
	output, err := s3Session.PutBucketWebsite(&s3.PutBucketWebsiteInput{
		Bucket: aws.String(config.DomainName),
		WebsiteConfiguration: &s3.WebsiteConfiguration{
			IndexDocument: &s3.IndexDocument{
				Suffix: aws.String("index.html"),
			},
			ErrorDocument: &s3.ErrorDocument{
				Key: aws.String("error.html"),
			},
		},
	})

	if err != nil {
		fmt.Printf("An error occured while making the website bucket: %s, %s", err, output)
	}

	log.Print("Bucket has been made into website")

	return nil
}

// DestroyBucket destroys the hosting bucket.
func DestroyBucket(config *StaticWebConfig, s3Session *s3.S3) error {

	_, err := s3Session.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(config.DomainName),
	})

	if err != nil {
		fmt.Printf("Unable to create bucket %q, %v", config.DomainName, err)
		return err
	}

	s3Session.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(config.DomainName),
	})

	log.Printf("Bucket %s has been destroyed", config.DomainName)

	return nil
}
