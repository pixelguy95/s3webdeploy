package deploy

import (
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
	CreateBucket(config, s3Session)
	SetBucketPermissions(config, s3Session)
	CreateBucketWebsite(config, s3Session)
	UploadWebFolder(config, sess)
	CreateCNameRecord(config, route53Session, &AliasConfig{DNSName: "blah"})
	return nil
}

func Cleanup(config *StaticWebConfig) error {

	sess, err := GetNewSession(config)

	if err != nil {
		return err
	}

	s3Session := s3.New(sess)
	DestroyBucket(config, s3Session)

	return nil
}
