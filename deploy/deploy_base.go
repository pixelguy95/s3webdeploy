package deploy

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Setup constucts a brand new s3-bucket website from scratch.
// Setup involves the following steps:
// 1. Creating the bucket
// 2. Setting the bucket reading permissions to public
// 3. Configure the bucket to act like a website.
// 4. Uploading the web-folder to the bucket.
//		The files needs the appropriate keys and the right mime-file-types to
//		act correctly on fetch.
// 5. Creates a new CNAME record pointing to the new bucket website.
//		Configured with the subdomain given in the config
// TODO: Move alias config into the cname record function
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

	url, region, _ := ExtractBucketWebsiteURL(config, s3Session)

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

// Cleanup  removes everything created in the Setup function. Basically a
// full reset to the pre-creation stage.
func Cleanup(config *StaticWebConfig) error {

	sess, err := GetNewSession(config)

	if err != nil {
		return err
	}

	s3Session := s3.New(sess)
	route53Session := route53.New(sess)

	url, region, _ := ExtractBucketWebsiteURL(config, s3Session)

	DeleteCNameRecord(config, route53Session, &AliasConfig{DNSName: *url, Region: *region})
	DestroyBucket(config, s3Session)

	return nil
}

// Update updates all the files in the bucket, simply uploads and overwrites
// everything in the bucket with the files in the local folder.
func Update(config *StaticWebConfig) error {
	sess, err := GetNewSession(config)
	if err != nil {
		return err
	}

	UploadWebFolder(config, sess)

	return nil
}
