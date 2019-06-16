package tests

import (
	"path/filepath"
	"testing"

	"github.com/pixelguy95/deploy-static-web/deploy"
)

func TestLoadConfigFile(t *testing.T) {

	path := filepath.Join("testdata", "config.json")
	t.Log(path)
	config := deploy.LoadConfigurations(path)

	if config.CredentialsName != "default" {
		t.Error("Error reading credentials")
	}

	if config.DomainName != "www.example.com" {
		t.Error("Error loading domain name")
	}

	if config.WebFolder != "./web-test" {
		t.Error("Error loading folder name")
	}
}
