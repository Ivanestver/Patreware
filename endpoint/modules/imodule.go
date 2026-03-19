package modules

import (
	"log"
)

type InfectionState = int

const (
	INFECTION_STATE_INFECTED InfectionState = iota
	INFECTION_STATE_CLEAN
	INFECTION_STATE_UNDEFINED
)

type Severity = int

const (
	SEVERITY_LOW Severity = iota
	SEVERITY_MEDIUM
	SEVERITY_HIGH
)

func SeverityToString(severity Severity) string {
	switch severity {
	case SEVERITY_LOW:
		return "Low"
	case SEVERITY_MEDIUM:
		return "Medium"
	case SEVERITY_HIGH:
		return "High"
	default:
		return "Undefined"
	}
}

type CheckProgress struct {
	PercentCompleted int
}

type CheckResult struct {
	AnalysisType string
	Path         string
	Severity     Severity
	Result       InfectionState
	ThreatName   string
}

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
	IsLoaded() bool
	IsSafe(path string, progressChan chan CheckProgress, resultChan chan CheckResult, errorChan chan error)
}

type BaseModule struct {
	Logger *log.Logger
}

func NewBaseModule(logger *log.Logger) BaseModule {
	return BaseModule{
		Logger: logger,
	}
}
