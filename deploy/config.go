package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
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
func LoadConfigurations(fileName string) *StaticWebConfig {
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	data := StaticWebConfig{}
	err = json.Unmarshal([]byte(file), &data)

	fmt.Println(string([]byte(file)))

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &data
}

func (conf *StaticWebConfig) SanityCheck() error {
	sess, err := GetNewSession(conf)

	if err != nil {
		return err
	}

	//s3Session := s3.New(sess)
	route53Session := route53.New(sess)

	output, err := route53Session.GetHostedZone(&route53.GetHostedZoneInput{
		Id: aws.String(conf.HostedZoneID),
	})

	if err != nil {
		return err
	}

	if conf.Subdomain != "" && !strings.HasSuffix(conf.Subdomain, ".") {
		conf.Subdomain = conf.Subdomain + "."
	}

	if strings.HasPrefix(conf.Subdomain, "www.") {
		conf.Subdomain = strings.TrimPrefix(conf.Subdomain, "www.")
	}

	conf.BucketName = "www." + conf.Subdomain + strings.TrimSuffix(*output.HostedZone.Name, ".")

	return nil
}
