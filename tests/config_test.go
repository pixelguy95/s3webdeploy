package tests

import (
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	config := LoadConfigurations("./config.json")

}
