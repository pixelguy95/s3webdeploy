package tests

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Tests the connection to the aws servers
// TODO: Fix so that it loads configurations from config file
func TestAwsConnection(t *testing.T) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})

	if err != nil {
		t.Error("Credentials error")
	}

	s3Session := s3.New(sess)

	_, err = s3Session.ListBuckets(nil)
	if err != nil {
		t.Error("Could not list buckets")
	}
}
