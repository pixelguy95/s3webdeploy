package deploy

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Setup(config *StaticWebConfig) error {

	sess, err := GetNewSession(config)

	if err != nil {
		return err
	}

	s3Session := s3.New(sess)
	route53Session := route53.New(sess)

	fmt.Println("Creating bucket")
	err = CreateBucket(config, s3Session)
	if err != nil {
		return err
	}

	fmt.Println("Setting bucket permissions")
	err = SetBucketPermissions(config, s3Session)
	if err != nil {
		return err
	}

	fmt.Println("Creating bucket website")
	err = CreateBucketWebsite(config, s3Session)
	if err != nil {
		return err
	}

	url, region, _ := ExtractBucketWebsiteUrl(config, s3Session)

	fmt.Println("Bucket website now operating @")
	fmt.Println(*url)
	fmt.Printf("In region %s, with the hosted zone id %s\n", *region, S3BucketHostedZoneMap[*region])

	err = UploadWebFolder(config, sess)
	if err != nil {
		return err
	}

	CreateCNameRecord(config, route53Session, &AliasConfig{DNSName: *url, Region: *region})

	return nil
}

func Cleanup(config *StaticWebConfig) error {

	sess, err := GetNewSession(config)

	if err != nil {
		return err
	}

	s3Session := s3.New(sess)
	route53Session := route53.New(sess)

	url, region, _ := ExtractBucketWebsiteUrl(config, s3Session)

	DeleteCNameRecord(config, route53Session, &AliasConfig{DNSName: *url, Region: *region})
	DestroyBucket(config, s3Session)

	return nil
}

func Update(config *StaticWebConfig) error {
	sess, err := GetNewSession(config)
	if err != nil {
		return err
	}

	UploadWebFolder(config, sess)

	return nil
}
