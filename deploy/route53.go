package deploy

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

// AliasConfig contains all the data needed to create a CNAME record pointing
// to a S3-bucket
type AliasConfig struct {
	DNSName string
	Region  string
}

// S3BucketHostedZoneMap is a map of region to hostedzone id
// Taken from the official documentation:
// https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region
var S3BucketHostedZoneMap = map[string]string{

	//United states
	"us-east-2": "Z2O1EMRO9K5GLX",
	"us-east-1": "Z3AQBSTGFYJSTF",
	"us-west-1": "Z2F56UZL2M1ACD",
	"us-west-2": "Z3BJ6K6RIION7M",

	//Asia pacific
	"ap-east-1":      "ZNB98KWMFR0R6",
	"ap-south-1":     "Z11RGJOFQNVJUP",
	"ap-northeast-3": "Z2YQB5RD63NC85",
	"ap-northeast-2": "Z3W03O7B5YMIYP",
	"ap-southeast-1": "Z3O0J2DXBE1FTB",
	"ap-southeast-2": "Z1WCIGYICN2BYD",
	"ap-northeast-1": "Z2M4EHUR26P7ZW",

	//Canada
	"ca-central-1": "Z1QDHH18159H29",

	//EU
	"eu-central-1": "Z21DNDUVLTQW6Q",
	"eu-west-1":    "Z1BKCTXD74EZPE",
	"eu-west-2":    "Z3GKZC51ZF0DB4",
	"eu-west-3":    "Z3R1K369G5AVDG",
	"eu-north-1":   "Z3BAZG2TWCNX0D",

	//South america
	"sa-east-1": "Z7KQH4QJS55SO",
}

func CreateCNameRecord(conf *StaticWebConfig, route53Session *route53.Route53, aliasConfig *AliasConfig) {

	fmt.Printf("Creating CNAME record with the name %s\n", conf.BucketName)
	fmt.Printf("Bucket in region %s, thus given hosted zone id %s\n", conf.Region, S3BucketHostedZoneMap[aliasConfig.Region])
	fmt.Printf("CNAME alias target %s\n", aliasConfig.DNSName)

	change := &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{

			Name: aws.String(conf.BucketName),
			Type: aws.String("CNAME"),

			AliasTarget: &route53.AliasTarget{
				EvaluateTargetHealth: aws.Bool(false),
				DNSName:              aws.String(aliasConfig.DNSName),
				HostedZoneId:         aws.String(S3BucketHostedZoneMap[aliasConfig.Region]),
			},
		},
		Action: aws.String(route53.ChangeActionCreate),
	}

	_, err := route53Session.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{change},
		},
		HostedZoneId: aws.String(conf.HostedZoneID),
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Created cname record")
}

func DeleteCNameRecord(conf *StaticWebConfig, route53Session *route53.Route53, aliasConfig *AliasConfig) {

	fmt.Printf("%v\n%v\n", conf, aliasConfig)
	change := &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{

			Name: aws.String(conf.BucketName),
			Type: aws.String("CNAME"),

			AliasTarget: &route53.AliasTarget{
				EvaluateTargetHealth: aws.Bool(false),
				DNSName:              aws.String(aliasConfig.DNSName),
				HostedZoneId:         aws.String(S3BucketHostedZoneMap[aliasConfig.Region]),
			},
		},
		Action: aws.String(route53.ChangeActionDelete),
	}

	_, err := route53Session.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{change},
		},
		HostedZoneId: aws.String(conf.HostedZoneID),
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Deleted cname record")
}
