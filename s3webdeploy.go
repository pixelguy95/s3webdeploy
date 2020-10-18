package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pixelguy95/s3webdeploy/deploy"
)

// Entry point of the program
func main() {

	if len(os.Args) < 2 {
		fmt.Println("No argument given")
		fmt.Println("Usage: s3webdeploy [opt] [(optional) config]")
		fmt.Println("")
		fmt.Println("possible options are: create, update, delete")
		return
	}

	option := strings.ToLower(os.Args[1])

	if option != "create" && option != "delete" && option != "update" {
		fmt.Printf("Unknown option %s, accepted options are create, update, delete\n", option)
		return
	}

	conf, err := handleConfigFile(option)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if option == deploy.CREATE {
		err = deploy.Setup(conf)
		if err != nil {
			if strings.HasPrefix(err.Error(), "BucketAlreadyOwnedByYou") {
				fmt.Println("You have already created this bucket")
			}

			fmt.Printf("%v\n", err)
		}

		fmt.Println("SUCCSESS!")
		fmt.Println("DNS records might take a few minutes to update")

	} else if option == deploy.DELETE {
		deploy.Cleanup(conf)
		if err != nil {
			fmt.Println("SUCCSESS!")
			fmt.Println("DNS records might take a few minutes to update")
		}

	} else if option == deploy.UPDATE {
		fmt.Println("Overwriting files in s3-bucket with new files")
		deploy.Update(conf)
		if err != nil {
			fmt.Println("SUCCSESS!")
			fmt.Println("All files overwritten")
		}
	}

}

// handleConfigFile loads the correct configuration file and parses it.
// in case of no config file found or malformed config file throws
// appropriate error.
// Also does a sanity check on the given configurations to see if they make
// any sense. If not throw nice error
// Otherwise returns pointer to new config struct
func handleConfigFile(option string) (*deploy.StaticWebConfig, error) {
	confFile := "./config.json"
	if len(os.Args) == 3 {
		confFile = os.Args[2]
	} else {
		fmt.Println("No config file given, using default value './config.json")
	}

	conf, err := deploy.LoadConfigurations(confFile)

	if err != nil {
		return nil, err
	}

	err = conf.SanityCheck(option)
	if err != nil {
		return nil, err
	}

	return conf, err
}
