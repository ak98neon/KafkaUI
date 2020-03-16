package web

import (
	"../kafka"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Error struct {
	ErrorMessage string
}

var isConfigured = false

func ConfigHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isConfigured {
			log.Println("Service isn't configured")
			http.Redirect(w, r, "/config", http.StatusSeeOther)
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}

func CommandPage(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	errorMsg := query.Get("error")

	t, err := template.ParseFiles("./src/page/main.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(writer, Error{ErrorMessage: errorMsg})
	if err != nil {
		log.Println("template executing error:", err)
	}
}

func ProduceFile(writer http.ResponseWriter, request *http.Request) {
	var errMsg string
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("Cannot parse multipart form: ", err)
		http.Redirect(writer, request, "/?error="+errMsg, http.StatusSeeOther)
		return
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		errMsg += "No such file"
		log.Println(errMsg)
		http.Redirect(writer, request, "/?error="+errMsg, http.StatusSeeOther)
		return
	}

	if file != nil {
		defer file.Close()
	}

	err = request.ParseForm()
	checkError(err)

	countValue := request.PostFormValue("count")
	parseInt, err := strconv.ParseInt(countValue, 10, 64)
	if err != nil || parseInt <= 0 {
		parseInt = 1
	}

	topic := request.PostFormValue("topic")
	if len(topic) <= 0 {
		errMsg += "Topic cannot be less or equal than 0"
		log.Println(errMsg)
		http.Redirect(writer, request, "/?error="+errMsg, http.StatusSeeOther)
		return
	}

	kafka.ProduceMessage(file, int(parseInt), topic)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func ConfigurationPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./src/page/configuration.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println("template executing error:", err)
	}
}

func Configuration(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	checkError(err)

	brokerServers := request.PostFormValue("broker_servers")
	kafka.BrokerList = strings.Split(brokerServers, ",")

	clientId := request.PostFormValue("clientId")
	kafka.ClientId = clientId

	isConfigured = true
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
