package main

import (
	"github.com/ak98neon/KafkaUI/src/web"
	"log"
	"net/http"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/config", web.ConfigurationPage)
	http.HandleFunc("/config/parse", web.Configuration)

	commandPage := http.HandlerFunc(web.CommandPage)
	producedFilePage := http.HandlerFunc(web.ProduceMessage)

	http.HandleFunc("/", web.ConfigHandler(commandPage))
	http.HandleFunc("/produce/file", web.ConfigHandler(producedFilePage))

	//Rest
	restKafkaInfo := http.HandlerFunc(web.GetKafkaInfo)
	http.HandleFunc("/kafka/info", web.ConfigHandlerRest(restKafkaInfo))
	log.Fatal(http.ListenAndServe(":9110", nil))
}
