package web

import (
	"github.com/ak98neon/KafkaUI/src/kafka"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Error struct {
	ErrorMessage string
}

type Response struct {
	KafkaInfo kafka.Info
	Error     Error
}

var IsConfigured = false

func ConfigHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsConfigured {
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

	info := kafka.GetKafkaInfo()
	err = t.Execute(writer, Response{KafkaInfo: info, Error: Error{ErrorMessage: errorMsg}})
	if err != nil {
		log.Println("template executing error:", err)
	}
}

func ProduceMessage(writer http.ResponseWriter, request *http.Request) {
	var errMsg string
	var isFile = true
	var inputMessage string

	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("Cannot parse multipart form: ", err)
		http.Redirect(writer, request, "/?error="+errMsg, http.StatusSeeOther)
		return
	}

	err = request.ParseForm()
	checkError(err)

	file, _, err := request.FormFile("file")
	if err != nil {
		isFile = false
		inputMessage = request.PostFormValue("input-message")
		if inputMessage == "" {
			errMsg += "No such file or string message"
			log.Println(errMsg)
			http.Redirect(writer, request, "/?error="+errMsg, http.StatusSeeOther)
		}
	}

	if file != nil {
		defer file.Close()
	}

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

	if isFile {
		osFile := MultipartFileToOsFile(file)
		fromFile := kafka.PrepareMessageFromFile(osFile, topic)
		kafka.ProduceMessage(fromFile, int(parseInt), topic)
	} else {
		fromString := kafka.PrepareMessageFromString(inputMessage, topic)
		kafka.ProduceMessage(fromString, int(parseInt), topic)
	}
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func ConfigurationPage(w http.ResponseWriter, _ *http.Request) {
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

	brokerServers := request.FormValue("broker_servers")
	kafka.BrokerList = strings.Split(brokerServers, ",")

	clientId := request.FormValue("clientId")
	kafka.ClientId = clientId

	if brokerServers == "" || clientId == "" {
		http.Redirect(writer, request, "/config", http.StatusSeeOther)
		return
	}

	IsConfigured = true
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
