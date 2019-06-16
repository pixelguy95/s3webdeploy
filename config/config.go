package config

import (
	"encoding/json"
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
	file, _ := ioutil.ReadFile(fileName)
	data := StaticWebConfig{}
	_ = json.Unmarshal([]byte(file), &data)
	return &data
}
