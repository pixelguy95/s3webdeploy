package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pixelguy95/deploy-static-web/deploy"
)

func main() {

	optPtr := flag.String("opt", "", "What the program should do, setup, update, cleanup")
	configPtr := flag.String("config", "./config.json", "Path of the config file, defaults to ./config.json")

	flag.Parse()

	if *optPtr == "" {
		log.Fatal("opt flags needs to be set, valid option are setup, update, and cleanup")
		fmt.Printf("%s\n", *configPtr)
		return
	}

	conf := deploy.LoadConfigurations("./config.json")
	fmt.Println(conf)

	deploy.Setup(conf)
	deploy.Cleanup(conf)
}
