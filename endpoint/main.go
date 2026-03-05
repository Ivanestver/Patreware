package main

import (
	"log"
	"os"
	"patrware-endpoint/config"
	"patrware-endpoint/modules"
	_ "patrware-endpoint/modules/hash_module"
	_ "patrware-endpoint/modules/signature_module"

	"github.com/hillu/go-yara/v4"
)

func main() {
	yara.NewCompiler()
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
		isInfected, err = currModule.IsSafe(os.Args[1])
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
