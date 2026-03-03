package main

import (
	"log"
	"os"
	"patrware-endpoint/config"
	"patrware-endpoint/modules/hash_module"
)

func main() {
	if len(os.Args) < 2 {
		panic("No file to check")
	}
	config.InitializeConfig()
	hashModule := hash_module.NewHashModule(log.Default())
	hashModule.LoadModule()
	isInfected, err := hashModule.Check(os.Args[1])
	if err != nil {
		panic(err.Error())
	} else if isInfected {
		log.Println("INFECTED!!!")
	} else {
		log.Println("NOT INFECTED")
	}
}
