package modules

import (
	"log"
)

var arrModules map[string]IModule = make(map[string]IModule)

func RegisterModule(module IModule) {
	if _, ok := arrModules[module.GetName()]; !ok {
		arrModules[module.GetName()] = module
	}
}

func GetAvailableModules() []string {
	availableModules := make([]string, len(arrModules))
	i := 0
	for name := range arrModules {
		availableModules[i] = name
		i++
	}
	return availableModules
}

func GetModuleByName(name string) IModule {
	if module, ok := arrModules[name]; ok {
		return module
	} else {
		return nil
	}
}

type IModule interface {
	GetName() string
	GetDescription() string
	LoadModule(args ...any) error
	IsSafe(path string) (bool, error)
}

type BaseModule struct {
	Logger *log.Logger
}

func NewBaseModule(logger *log.Logger) BaseModule {
	return BaseModule{
		Logger: logger,
	}
}

type CheckResult struct {
	ModuleName string
	Result     bool
}
