/*
Package handlers is blah-blah-blah
*/
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"patrware/server/hub"
)

func init() {
	webHub.Handlers["/"] = mainPageHandler
}

type _EndpointViewModel struct {
	IP string
}

func mainPageHandler(writer http.ResponseWriter, req *http.Request) {
	mainHTMLTemplatePath := getTemplatesPath("main.html")
	mainHTMLTemplate, err := template.ParseFiles(mainHTMLTemplatePath)
	if err != nil {
		panic(err.Error())
	}

	endpoints := hub.GetAllEndpoints()
	data := make([]_EndpointViewModel, len(endpoints))
	for idx, endpoint := range endpoints {
		conn, err := hub.GetConnectionAssociatedWithEndpoint(endpoint.GetID())
		if err != nil {
			log.Println(err.Error())
			continue
		}
		data[idx].IP = conn.SocketConn.LocalAddr().String()
	}
	if err = mainHTMLTemplate.Execute(writer, data); err != nil {
		log.Printf("%v\n", err)
	}
}
