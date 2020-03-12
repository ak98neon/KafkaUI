package main

import (
	"./src/web"
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", web.CommandPage)
	http.HandleFunc("/produce/file", web.ProduceFile)
	log.Fatal(http.ListenAndServe(":9100", nil))
}
