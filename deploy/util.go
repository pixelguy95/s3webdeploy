package deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetNewSession(config *StaticWebConfig) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", config.CredentialsName),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}
