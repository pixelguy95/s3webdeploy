package deploy

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	log.Print("Bucket reading permission set to public")

	return nil
}

func UploadWebFolder(config *StaticWebConfig, sess *session.Session) error {

	uploader := s3manager.NewUploader(sess)

	err := filepath.Walk(config.WebFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fileContent, err := os.Open(path)
			key := strings.Replace(path, "\\", "/", 100)

			if strings.HasPrefix(config.WebFolder, "./") {
				key = strings.Replace(key, strings.TrimSuffix("./", config.WebFolder)+"/", "", 1)
			} else {
				key = strings.Replace(key, config.WebFolder+"/", "", 1)
			}

			//TODO: Seperate file with mime defenitions
			extn := filepath.Ext(path)
			contentType := "application/octet-stream"
			if extn == ".html" {
				contentType = "text/html"
			} else if extn == ".pdf" {
				contentType = "application/pdf"
			} else if extn == ".css" {
				contentType = "text/css"
			} else if extn == ".js" {
				contentType = "application/javascript"
			} else if extn == ".png" || extn == ".jpg" || extn == ".gif" {
				contentType = "image/" + extn
			}

			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket:      aws.String(config.DomainName),
				Key:         aws.String(key),
				Body:        fileContent,
				ContentType: aws.String(contentType),
			})

			if err == nil {
				log.Printf("%s", key)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
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

func ExtractBucketWebsiteUrl(config *StaticWebConfig, s3Session *s3.S3) error {
	output, err := s3Session.GetBucketWebsite(&s3.GetBucketWebsiteInput{
		Bucket: aws.String(config.DomainName),
	})
	fmt.Println(output.String())
	fmt.Println(err)

	return nil
}

// DestroyBucket destroys the hosting bucket.
func DestroyBucket(config *StaticWebConfig, s3Session *s3.S3) error {

	list, err := s3Session.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(config.DomainName),
	})

	for _, l := range list.Contents {
		s3Session.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(config.DomainName),
			Key:    l.Key,
		})
	}

	_, err = s3Session.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(config.DomainName),
	})

	if err != nil {
		fmt.Printf("Unable to destroy bucket %q, %v", config.DomainName, err)
		return err
	}

	s3Session.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(config.DomainName),
	})

	log.Printf("Bucket %s has been destroyed", config.DomainName)

	return nil
}
