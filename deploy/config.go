package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
)

// StaticWebConfig defines the settings for the deployment
type StaticWebConfig struct {
	HostedZoneID    string `json:"hosted_zone_id"`
	Subdomain       string `json:"subdomain"`
	CredentialsName string `json:"credentials_name"`
	WebFolder       string `json:"folder"`
	BucketName      string
	Region          string `json:"region"`
}

// LoadConfigurations loads the configurations from the file given and
// returns the struct
func LoadConfigurations(fileName string) (*StaticWebConfig, error) {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Loading '%s' as config file:\n", fileName)

	data := StaticWebConfig{}
	err = json.Unmarshal([]byte(file), &data)

	fmt.Println(string([]byte(file)))

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// SanityCheck aChecks all the details given in the config file.
// TODO: Check if IAM user has the correct priviledges
func (conf *StaticWebConfig) SanityCheck(option string) error {

	fmt.Println("Performing sanity check and cleanup on configurations")
	sess, err := GetNewSession(conf)

	if err != nil {
		return err
	}

	s3Session := s3.New(sess)
	route53Session := route53.New(sess)

	output, err := route53Session.GetHostedZone(&route53.GetHostedZoneInput{
		Id: aws.String(conf.HostedZoneID),
	})

	if err != nil {
		if err.Error() == "SharedCredsLoad: failed to get profile" {
			fmt.Printf("Could not find credentials named '%s'\n", conf.CredentialsName)
		} else if strings.HasPrefix(err.Error(), "NoSuchHostedZone") {
			fmt.Printf("")
		}

		return err
	}

	if conf.Subdomain != "" && !strings.HasSuffix(conf.Subdomain, ".") {
		conf.Subdomain = conf.Subdomain + "."
	}

	conf.BucketName = conf.Subdomain + strings.TrimSuffix(*output.HostedZone.Name, ".")

	if option == CREATE {
		_, err := s3Session.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(conf.BucketName),
		})

		if err != nil && !strings.HasPrefix(err.Error(), "NotFound") {
			return err
		} else if err == nil {
			return errors.New("Bucket already exists, did you mean to update?")
		}

		out, err := route53Session.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(conf.HostedZoneID),
		})

		for _, record := range out.ResourceRecordSets {
			if *record.Name == conf.BucketName+"." {
				return errors.New("The record for " + conf.BucketName + " already exists, did you mean to delete?")
			}
		}
	}

	fmt.Printf("Clean configurations: \n%v\n", conf)

	return nil
}
