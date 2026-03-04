package main

import (
	"log"
	"os"
	"patrware-endpoint/config"
	"patrware-endpoint/modules"
	_ "patrware-endpoint/modules/hash_module"
)

func main() {
	if len(os.Args) < 2 {
		panic("No file to check")
	}
	config.InitializeConfig()
	availableModules := modules.GetAvailableModules()
	isInfected := false
	for _, moduleName := range availableModules {
		currModule := modules.GetModuleByName(moduleName)
		err := currModule.LoadModule()
		if err != nil {
			log.Println(err.Error())
			return
		}
		isInfected, err = currModule.Check(os.Args[1])
		if err != nil {
			panic(err.Error())
		} else if isInfected {
			break
		} else {
			continue
		}
	}
	if isInfected {
		log.Println("INFECTED!!!")
	} else {
		log.Println("NOT INFECTED")
	}
}
