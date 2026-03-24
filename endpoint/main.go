package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"patrware/endpoint/config"
	"patrware/endpoint/modules"
	_ "patrware/endpoint/modules/hash_module"
	_ "patrware/endpoint/modules/signature_module"
	pb "patrware/proto"
	"patrware/structs"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type Modules struct {
	availableModules []string
	modulesList      map[string]modules.IModule
}

func NewModules() *Modules {
	return &Modules{
		availableModules: modules.GetAvailableModules(),
		modulesList:      make(map[string]modules.IModule),
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
	checker := NewChecker(stream)
	checker.Check(req.Path)
	return nil
}

var modulesStorage *Modules

func main() {
	makeConnectionWithAdminServer()
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
	for _, moduleName := range modulesStorage.GetAvailableModules() {
		if module, err := modulesStorage.GetModule(moduleName); err == nil || !module.IsLoaded() {
			if err = module.LoadModule(); err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}

func mainLoop(listener net.Listener) {
	scanner := grpc.NewServer()
	pb.RegisterScannerServiceServer(scanner, &ScannerServer{})
	if err := scanner.Serve(listener); err != nil {
		fmt.Println(err.Error())
	}
}

type Checker struct {
	stream pb.ScannerService_StartScanServer
}

func NewChecker(stream pb.ScannerService_StartScanServer) *Checker {
	return &Checker{
		stream: stream,
	}
}

func (checker *Checker) Check(path string) {
	availableModules := modulesStorage.GetAvailableModules()
	files, err := checker.defineSetOfFilesToCheck(path)
	if err != nil {
		checker.stream.Send(&pb.ScanEvent{
			ErrorMsg: err.Error(),
		})
		return
	}
	for _, file := range files {
		checker.checkFile(file, availableModules)
	}
}

func (checker *Checker) checkFile(filepath string, availableModules []string) {
	progressChan := make(chan modules.CheckProgress)
	resultChan := make(chan modules.CheckResult)
	errorChan := make(chan error)
	var wg sync.WaitGroup
	for _, moduleName := range availableModules {
		currModule := modules.GetModuleByName(moduleName)
		if !currModule.IsLoaded() {
			if err := currModule.LoadModule(); err != nil {
				log.Println(err.Error())
				continue
			}
		}
		wg.Go(func() {
			currModule.IsSafe(filepath, progressChan, resultChan, errorChan)
		})
	}

	go func() {
		wg.Wait()
		close(progressChan)
		close(resultChan)
		close(errorChan)
	}()

	isRunning := true
	results := make([]modules.CheckResult, 0)
	for isRunning {
		select {
		case <-progressChan:
			// contribution := (float64(process.PercentCompleted) * 0.01) *
			// 	(1.0 / float64(len(availableModules))) * 100
			// currPercent += contribution
			// checker.stream.Send(&pb.ScanEvent{
			// 	ProgressPercent: int32(currPercent),
			// })
		case result, ok := <-resultChan:
			if ok {
				results = append(results, result)
			} else {
				isInfected := false
				for _, r := range results {
					isInfected = isInfected || r.Result == modules.INFECTION_STATE_INFECTED
				}
				checker.stream.Send(&pb.ScanEvent{
					CurrentFile:     filepath,
					ProgressPercent: 100,
					VirusFound:      isInfected,
					ThreatName:      "Some threat",
					ErrorMsg:        "",
				})
				isRunning = false
			}
		case err, ok := <-errorChan:
			if ok {
				checker.stream.Send(&pb.ScanEvent{
					CurrentFile:     filepath,
					ProgressPercent: 100,
					VirusFound:      false,
					ThreatName:      "Unknown",
					ErrorMsg:        err.Error(),
				})
				isRunning = false
			} else {

			}
		}
	}
}

func (checker *Checker) defineSetOfFilesToCheck(root string) ([]string, error) {
	if fileInfo, err := os.Stat(root); err != nil {
		return nil, err
	} else {
		if fileInfo.Mode().IsRegular() {
			return []string{root}, nil
		} else {
			files := make([]string, 0)
			if err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				files = append(files, path)
				return nil
			}); err != nil {
				return nil, err
			}
			return files, nil
		}
	}
}

func makeConnectionWithAdminServer() {
	serverURL := "ws://localhost:60000/ws"

	for {
		log.Println("Пытаюсь подключиться к серверу...")

		// 1. Устанавливаем соединение
		conn, resp, err := websocket.DefaultDialer.Dial(serverURL, nil)
		if err != nil {
			log.Printf("Ошибка подключения: %v. Повтор через 5 секунд...", err)
			if resp != nil {
				log.Printf("resp: %v", resp)
			}
			time.Sleep(5 * time.Second)
			continue
		}

		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("Ошибка при получении имени компьютера: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		info := structs.EndpointInfo{
			Name:          hostname,
			SecurityState: structs.SECURITY_STATE_CLEAN,
		}
		if err = conn.WriteJSON(&info); err != nil {
			log.Printf("Ошибка при отправке регистрационного сообщения, %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Подключено!")

		// 2. Обрабатываем соединение
		handleConnection(conn)

		// Если handleConnection завершился, значит связь прервалась
		log.Println("Связь потеряна.")
		time.Sleep(5 * time.Second)
	}
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	// Горутина для чтения ответов от сервера
	// go func() {
	// 	for {
	// 		var msg protocol.Message
	// 		err := conn.ReadJSON(&msg)
	// 		if err != nil {
	// 			log.Println("Ошибка чтения:", err)
	// 			return
	// 		}
	// 		log.Printf("Получено от сервера: %s", msg.Content)
	// 	}
	// }()

	// Основной цикл отправки данных (например, статус агента раз в 10 сек)
	for {
		// msg := protocol.Message{
		// 	From:    "Agent-007",
		// 	Content: "Я жив, всё ок!",
		// }

		// err := conn.WriteJSON(msg)
		// if err != nil {
		// 	log.Println("Ошибка отправки:", err)
		// 	return // Выходим, чтобы сработал реконнект в main
		// }

		time.Sleep(10 * time.Second)
	}
}
