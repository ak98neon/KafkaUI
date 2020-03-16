package main

import (
	"./src/web"
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/config", web.ConfigurationPage)
	http.HandleFunc("/config/parse", web.Configuration)

	commandPage := http.HandlerFunc(web.CommandPage)
	producedFilePage := http.HandlerFunc(web.ProduceFile)

	http.HandleFunc("/", web.ConfigHandler(commandPage))
	http.HandleFunc("/produce/file", web.ConfigHandler(producedFilePage))
	log.Fatal(http.ListenAndServe(":9100", nil))
}
