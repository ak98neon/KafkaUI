package main

import (
	"./src/web"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", web.CommandPage)
	http.HandleFunc("/produce/file", web.ProduceFile)
	log.Fatal(http.ListenAndServe(":9100", nil))

	//filePath := flag.String("file", "C:/Users/arkudrya/IdeaProjects/kafkaProducer/src/resources/routing_with_context.xml", "Path to the incoming file")
	////filePath := flag.String("file", "C:/Users/arkudrya/IdeaProjects/kafkaProducer/src/resources/high-overdue.in.xml", "Path to the incoming file")
	////filePath := flag.String("file", "C:/Users/arkudrya/IdeaProjects/kafkaProducer/src/resources/5_min.xml", "Path to the incoming file")
	//flag.Parse()
	//
	//kafka.ProduceMessage(filePath)
}
