package main

import (
	"context"
	"fmt"
	"io"
	"os"

	pb "patrware/proto"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UIScanEvent struct {
	// File     string `json:"file"`
	// Progress int    `json:"progress"`
	// IsVirus  bool   `json:"isVirus"`
	pb.ScanEvent
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetFilePathToScan() string {
	options := runtime.OpenDialogOptions{}
	options.DefaultDirectory, _ = os.UserHomeDir()
	options.Title = "Выбрать файл для сканирования"
	options.ShowHiddenFiles = true
	if path, err := runtime.OpenFileDialog(a.ctx, options); err == nil {
		return path
	} else {
		return ""
	}
}

func (a *App) GetDirPathToScan() string {
	options := runtime.OpenDialogOptions{}
	options.DefaultDirectory, _ = os.UserHomeDir()
	options.Title = "Выбрать директорию для сканирования"
	if path, err := runtime.OpenDirectoryDialog(a.ctx, options); err == nil {
		return path
	} else {
		return ""
	}
}

func (a *App) StartScan(path string) {
	conn, err := grpc.NewClient("localhost:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := pb.NewScannerServiceClient(conn)
	stream, err := client.StartScan(a.ctx, &pb.ScanRequest{Path: path})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go func() {
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				runtime.EventsEmit(a.ctx, "scan_complete", event)
				break
			}
			if err != nil {
				fmt.Printf("Ошибка стрима: %v", err)
				break
			}

			// 4. Отправляем данные во фронтенд через события Wails
			runtime.EventsEmit(a.ctx, "scan_progress", event)
		}
	}()
}

func (a *App) GenerateScanEvent(event UIScanEvent) {
}
