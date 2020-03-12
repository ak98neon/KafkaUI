package web

import (
	"../kafka"
	"html/template"
	"log"
	"net/http"
)

func CommandPage(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("./src/page/main.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(writer, nil)
	if err != nil {
		log.Println("template executing error:", err)
	}
}

func ProduceFile(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("Cannot parse multipart form: ", err)
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	kafka.ProduceMessage(file)

	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
