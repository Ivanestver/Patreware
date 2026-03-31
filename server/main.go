package main

import (
	"fmt"
	"net/http"
	"patrware/server/hub"
	"patrware/server/web/handlers"
)

func main() {
	hub.InitHub()
	handlers.SetupHandlers()

	if err := http.ListenAndServe(":60000", nil); err != nil {
		fmt.Println(err.Error())
	}
}
