package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// StaticWebConfig defines the settings for the deployment
type StaticWebConfig struct {
	DomainName      string `json:"domain"`
	CredentialsName string `json:"credentials_name"`
	WebFolder       string `json:"folder"`
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

	if err != nil {
		fmt.Println(err)
		return nil
	}

	data.DomainName = "www." + data.DomainName
	return &data
}
