package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	optPtr := flag.String("opt", "", "What the program should do, setup, update, cleanup")
	configPtr := flag.String("config", "./config.json", "Path of the config file, defaults to ./config.json")
	webFolderPtr := flag.String("folder", "", "Path of the folder you want to upload")

	flag.Parse()

	if *optPtr == "" || *webFolderPtr == "" {
		log.Fatal("opt and folder flags needs to be set")
		return
	}

	fmt.Printf("%s\n", *configPtr)
}
