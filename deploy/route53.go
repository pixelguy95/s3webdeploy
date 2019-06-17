package deploy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type AliasConfig struct {
	DNSName string
}

func extractHostedZoneID(name string, route53Session *route53.Route53) (*string, error) {
	zones, err := route53Session.ListHostedZones(nil)

	if err != nil {
		fmt.Print(err)
	}

	routeStyleDomainName := strings.TrimPrefix(name, "www.") + "."

	for _, i := range zones.HostedZones {
		if *i.Name == routeStyleDomainName {
			return i.Id, nil
		}
	}

	return nil, errors.New("Could not find hosted zone with that name")
}

func CreateCNameRecord(conf *StaticWebConfig, route53Session *route53.Route53, aliasConfig *AliasConfig) {

	hostedZoneId, err := extractHostedZoneID(conf.DomainName, route53Session)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*hostedZoneId)

	// change := &route53.Change{
	// 	ResourceRecordSet: &route53.ResourceRecordSet{
	//
	// 		AliasTarget: &route53.AliasTarget{
	// 			EvaluateTargetHealth: aws.Bool(false),
	// 			DNSName:              aws.String(aliasConfig.DNSName),
	// 		},
	// 	},
	// }
	//
	// route53Session.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
	// 	ChangeBatch: &route53.ChangeBatch{
	// 		Changes: []*route53.Change{change},
	// 	},
	//
	// })
}
