package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pixelguy95/deploy-static-web/deploy"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No argument given")
		fmt.Println("Usage: s3webdeploy [opt] [(optional) config]")
		fmt.Println("")
		fmt.Println("possible options are: create, update, delete")
		return
	}

	fmt.Println("First option: " + os.Args[1])
	fmt.Println("Second option: " + os.Args[2])

	if os.Args[1] == "create" {

	}

	conf := deploy.LoadConfigurations("./config.json")

	if conf == nil {
		return
	}

	err := conf.SanityCheck()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(conf)

	err = deploy.Setup(conf)
	if err != nil {
		if strings.HasPrefix(err.Error(), "BucketAlreadyOwnedByYou") {
			fmt.Println("You have already created this bucket")
		}

		fmt.Printf("%v\n", err)
	}
	//deploy.Cleanup(conf)
}
