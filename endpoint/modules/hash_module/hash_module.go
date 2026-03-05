package hash_module

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"log"
	"os"
	"patrware-endpoint/config"
	"patrware-endpoint/modules"
	"slices"
)

func init() {
	modules.RegisterModule(NewHashModule(log.Default()))
}

type HashValue = string

type HashModule struct {
	modules.BaseModule
	md5Hashes    []string
	sha256Hashes []string
}

func NewHashModule(logger *log.Logger) *HashModule {
	return &HashModule{
		BaseModule:   modules.NewBaseModule(logger),
		md5Hashes:    make([]string, 0),
		sha256Hashes: make([]string, 0),
	}
}

func (module *HashModule) GetName() string {
	return "Hash Module"
}

func (module *HashModule) GetDescription() string {
	return "The module checks a file taking its hash value and comparing with the known ones"
}

func (module *HashModule) LoadModule(args ...any) error {
	conf := config.GetConfig()
	md5_chan := make(chan error)
	sha256_chan := make(chan error)
	go func(c chan error) {
		err := module.loadHashes(conf.Hashes.MD5HashPath, &module.md5Hashes)
		c <- err
	}(md5_chan)
	go func(c chan error) {
		err := module.loadHashes(conf.Hashes.SHA256HashPath, &module.sha256Hashes)
		c <- err
	}(sha256_chan)
	return errors.Join(<-md5_chan, <-sha256_chan)
}

func (module *HashModule) loadHashes(hashpath string, hashStorage *[]string) error {
	entries, err := os.ReadDir(hashpath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file, err := os.Open(hashpath + entry.Name())
		if err != nil {
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			*hashStorage = append(*hashStorage, scanner.Text())
		}
	}
	return nil
}

func (module *HashModule) IsSafe(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	md5Result, _, err := module.checkMD5Hash(file)
	if err != nil {
		return false, nil
	}
	sha256Result, _, err := module.checkSHA256Hash(file)
	if err != nil {
		return false, nil
	}
	return md5Result || sha256Result, nil
}

func (module *HashModule) checkMD5Hash(file *os.File) (bool, HashValue, error) {
	md5func := md5.New()
	if hashValue, err := module.calcHash(&md5func, file); err != nil {
		return false, "", err
	} else {
		return module.checkHash(hashValue, module.md5Hashes), hashValue, nil
	}
}

func (module *HashModule) checkSHA256Hash(file *os.File) (bool, HashValue, error) {
	sha256func := sha256.New()
	if hashValue, err := module.calcHash(&sha256func, file); err != nil {
		return false, "", err
	} else {
		return module.checkHash(hashValue, module.md5Hashes), hashValue, nil
	}
}

func (module *HashModule) calcHash(hash *hash.Hash, file *os.File) (HashValue, error) {
	if _, err := io.Copy(*hash, file); err != nil {
		return "", err
	}
	fileHash := hex.EncodeToString((*hash).Sum(nil))
	return fileHash, nil
}

func (module *HashModule) checkHash(hashValue HashValue, hashes []string) bool {
	hashIndex := slices.Index(hashes, hashValue)
	return hashIndex != -1
}
