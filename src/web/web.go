package web

import (
	"../kafka"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	checkError(err)

	defer file.Close()

	err = request.ParseForm()
	checkError(err)

	countValue := request.PostFormValue("count")
	parseInt, err := strconv.ParseInt(countValue, 10, 64)
	if err != nil || parseInt <= 0 {
		parseInt = 1
	}

	kafka.ProduceMessage(file, int(parseInt))
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
