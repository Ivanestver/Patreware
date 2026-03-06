package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"patrware-endpoint/config"
	"patrware-endpoint/modules"
	_ "patrware-endpoint/modules/hash_module"
	_ "patrware-endpoint/modules/signature_module"
	"slices"
)

type Modules struct {
	availableModules []string
	modulesList      map[string]modules.IModule
}

func NewModules() *Modules {
	return &Modules{
		availableModules: modules.GetAvailableModules(),
	}
}

func (mods *Modules) GetAvailableModules() []string {
	return mods.availableModules
}

func (mods *Modules) GetModule(moduleName string) (modules.IModule, error) {
	if idx := slices.Index(mods.availableModules, moduleName); idx == -1 {
		return nil, errors.New("No module named " + moduleName)
	}
	if module, ok := mods.modulesList[moduleName]; ok {
		return module, nil
	} else {
		m := modules.GetModuleByName(moduleName)
		if err := m.LoadModule(); err != nil {
			return nil, err
		}
		mods.modulesList[moduleName] = m
		return m, nil
	}
}

var modulesStorage *Modules

func main() {
	host := "localhost"
	port := 50000
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err.Error())
	}
	defer listener.Close()
	fmt.Printf("[INFO] Server started listening to the port %d\n", port)
	config.InitializeConfig()
	modulesStorage = NewModules()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)
	defer conn.Close()
	var req CheckRequest
	resp := CheckResponse{}
	if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
		resp.Ok = false
		resp.Error = err.Error()
		sendResponce(&resp, writer)
		return
	}
	isInfected := false
	availableModules := modulesStorage.GetAvailableModules()
	for _, moduleName := range availableModules {
		currModule := modules.GetModuleByName(moduleName)
		err := currModule.LoadModule()
		if err != nil {
			log.Println(err.Error())
			return
		}
		isInfected, err = currModule.IsSafe(os.Args[1])
		if err != nil {
			log.Println(err.Error())
		} else if isInfected {
			break
		} else {
			continue
		}
	}
	if isInfected {
		resp.Ok = true
		resp.Result = "INFECTED!!!"
	} else {
		resp.Ok = false
		resp.Result = "NOT INFECTED!!!"
	}

	sendResponce(&resp, writer)
}

func sendResponce(resp *CheckResponse, writer *bufio.Writer) {
	if bytes, err := json.Marshal(*resp); err == nil {
		writer.Write(bytes)
	} else {
		bytes, _ := json.Marshal(CheckResponse{
			Ok:    false,
			Error: "Internal Server Error",
		})
		writer.Write(bytes)
	}
}

type CheckRequest struct {
	Path string `json:"path"`
}

type CheckResponse struct {
	Ok     bool   `json:"ok"`
	Result string `json:"result"`
	Error  string `json:"error"`
}
