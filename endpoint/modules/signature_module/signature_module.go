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

type SignatureModule struct {
	modules.BaseModule
	Rules *yara.Rules
}

func NewSignatureModule(logger *log.Logger) *SignatureModule {
	return &SignatureModule{}
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
	if rules, err := compiler.GetRules(); err != nil {
		return err
	} else {
		module.Rules = rules
		return nil
	}
}

func (module *SignatureModule) IsSafe(path string) (bool, error) {
	var matches yara.MatchRules
	file, err := os.Open(path)
	if err != nil {
		return true, err
	}
	module.Rules.ScanFileDescriptor(file.Fd(), 0, 0, &matches)
	if len(matches) == 0 {
		return true, nil
	}
	return false, nil
}
