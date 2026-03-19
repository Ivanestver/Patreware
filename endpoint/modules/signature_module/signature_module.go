package signature_module

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"patrware-endpoint/config"
	"patrware-endpoint/modules"

	yara "github.com/hillu/go-yara/v4"
)

func init() {
	modules.RegisterModule(NewSignatureModule(log.Default()))
}

type _CustomScanCallback struct {
	Matches        yara.MatchRules
	ProgressChan   chan modules.CheckProgress
	RulesCount     float64
	currRulesCount float64
}

func (cb *_CustomScanCallback) RuleMatching(sc *yara.ScanContext, r *yara.Rule) (abort bool, err error) {
	abort, err = cb.Matches.RuleMatching(sc, r)
	cb.currRulesCount++
	cb.ProgressChan <- modules.CheckProgress{
		PercentCompleted: int(cb.currRulesCount / cb.RulesCount * 100.0),
	}
	return
}

type SignatureModule struct {
	modules.BaseModule
	Rules    *yara.Rules
	isLoaded bool
}

func NewSignatureModule(logger *log.Logger) *SignatureModule {
	return &SignatureModule{isLoaded: false}
}

func (module *SignatureModule) GetName() string {
	return "Behavior Module"
}

func (module *SignatureModule) GetDescription() string {
	return "This module looks at the file parts to find malicious code"
}

func (module *SignatureModule) LoadModule(args ...any) error {
	compiler, err := yara.NewCompiler()
	if err != nil {
		return nil
	}
	conf := config.GetConfig()
	if err = filepath.Walk(conf.Signatures.Path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".yar" || filepath.Ext(path) != ".yara" {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		return compiler.AddFile(file, path)
	}); err != nil {
		return err
	}
	module.isLoaded = true
	if rules, err := compiler.GetRules(); err != nil {
		return err
	} else {
		module.Rules = rules
		return nil
	}
}

func (module *SignatureModule) IsLoaded() bool {
	return module.isLoaded
}

func (module *SignatureModule) IsSafe(path string, progressChan chan modules.CheckProgress, resultChan chan modules.CheckResult, errorChan chan error) {
	checkResult := modules.CheckResult{
		AnalysisType: module.getAnalysisType(),
		Path:         path,
		Severity:     modules.SEVERITY_HIGH,
		Result:       modules.INFECTION_STATE_UNDEFINED,
	}
	file, err := os.Open(path)
	if err != nil {
		errorChan <- err
		return
	}
	callback := _CustomScanCallback{
		ProgressChan: progressChan,
		RulesCount:   float64(len(module.Rules.GetRules())),
	}
	ruleMatches := yara.MatchRules{}
	module.Rules.ScanFileDescriptor(file.Fd(), 0, 0, &ruleMatches)
	if len(callback.Matches) == 0 {
		checkResult.Result = modules.INFECTION_STATE_CLEAN
	} else {
		checkResult.ThreatName = callback.Matches[0].Metas[0].Identifier
		checkResult.Result = modules.INFECTION_STATE_INFECTED
	}
	resultChan <- checkResult
}

func (module *SignatureModule) getAnalysisType() string {
	return "Signature Analysis"
}
