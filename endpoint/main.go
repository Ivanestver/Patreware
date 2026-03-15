package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"patrware-endpoint/config"
	"patrware-endpoint/modules"
	_ "patrware-endpoint/modules/hash_module"
	_ "patrware-endpoint/modules/signature_module"
	pb "patrware/proto"
	"slices"

	"google.golang.org/grpc"
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

type ScannerServer struct {
	pb.UnimplementedScannerServiceServer
}

func (s *ScannerServer) StartScan(req *pb.ScanRequest, stream pb.ScannerService_StartScanServer) error {
	isInfected := checkIfInfected(req.Path)
	stream.Send(&pb.ScanEvent{
		CurrentFile:     req.Path,
		ProgressPercent: 100,
		VirusFound:      isInfected,
		ThreatName:      "Some motherfucker",
	})
	return nil
}

var modulesStorage *Modules

func main() {
	listener := makeListener()
	defer listener.Close()
	configure()
	mainLoop(listener)
}

func makeListener() net.Listener {
	host := "localhost"
	port := 50000
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[INFO] Server started listening to the port %d\n", port)
	return listener
}

func configure() {
	config.InitializeConfig()
	modulesStorage = NewModules()
}

func mainLoop(listener net.Listener) {
	scanner := grpc.NewServer()
	pb.RegisterScannerServiceServer(scanner, &ScannerServer{})
	if err := scanner.Serve(listener); err != nil {
		fmt.Println(err.Error())
	}
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		continue
	// 	}
	// 	go handleConnection(conn)
	// }
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	req := getRequest(conn)
	if req == nil {
		return
	}
	isInfected := checkIfInfected(req.Path)
	resp := makeResponseIfInfected(isInfected)
	sendResponce(resp, conn)
}

func getRequest(conn net.Conn) *CheckRequest {
	scanner := bufio.NewScanner(conn)
	req := &CheckRequest{}
	if err := json.Unmarshal(scanner.Bytes(), req); err != nil {
		resp := CheckResponse{}
		resp.Ok = false
		resp.Error = err.Error()
		sendResponce(&resp, conn)
		return nil
	}
	return req
}

func sendResponce(resp *CheckResponse, conn net.Conn) {
	writer := bufio.NewWriter(conn)
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

func checkIfInfected(path string) bool {
	availableModules := modulesStorage.GetAvailableModules()
	for _, moduleName := range availableModules {
		currModule := modules.GetModuleByName(moduleName)
		err := currModule.LoadModule()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		isInfected, err := currModule.IsSafe(path)
		if err != nil {
			log.Println(err.Error())
		} else if isInfected {
			return true
		} else {
			continue
		}
	}
	return false
}

func makeResponseIfInfected(isInfected bool) *CheckResponse {
	resp := &CheckResponse{}
	if isInfected {
		resp.Ok = true
		resp.Result = "INFECTED!!!"
	} else {
		resp.Ok = false
		resp.Result = "NOT INFECTED!!!"
	}
	return resp
}

type CheckRequest struct {
	Path string `json:"path"`
}

type CheckResponse struct {
	Ok     bool   `json:"ok"`
	Result string `json:"result"`
	Error  string `json:"error"`
}
