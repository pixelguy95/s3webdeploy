package deploy

import (
	"log"

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
	err = CreateBucket(config, s3Session)
	if err != nil {
		return err
	}

	err = SetBucketPermissions(config, s3Session)
	if err != nil {
		return err
	}

	err = CreateBucketWebsite(config, s3Session)
	if err != nil {
		return err
	}

	url, region, _ := ExtractBucketWebsiteUrl(config, s3Session)

	log.Println(*url, *region, S3BucketHostedZoneMap[*region])

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
