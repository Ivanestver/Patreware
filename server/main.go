package main

import (
	"fmt"
	"net/http"
	"patrware/server/hub"
)

func main() {
	hubInstance := hub.NewHub()
	http.HandleFunc("/ws", hubInstance.ServeConnection)

	if err := http.ListenAndServe(":60000", nil); err != nil {
		fmt.Println(err.Error())
	}
}
