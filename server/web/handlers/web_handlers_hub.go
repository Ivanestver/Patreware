package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func getTemplatesPath(htmlFileName string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	return filepath.Join(wd, "web", "templates", htmlFileName)
}

type _WebHandler = http.HandlerFunc
type _Path = string

type _WebHandlersHub struct {
	Handlers map[_Path]_WebHandler
}

var webHub _WebHandlersHub = func() _WebHandlersHub {
	h := _WebHandlersHub{
		Handlers: make(map[_Path]_WebHandler),
	}
	h.Handlers["/"] = mainPageHandler
	return h
}()

func SetupHandlers() {
	for path, handler := range webHub.Handlers {
		http.HandleFunc(path, handler)
	}
	staticHandler := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
}
