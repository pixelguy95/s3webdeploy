package deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// GetNewSession takes in a config and returns the proper aws sdk session
// TODO: handle s3Session and route53Session creation too?
func GetNewSession(config *StaticWebConfig) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewSharedCredentials("", config.CredentialsName),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}
