package web

import (
	"encoding/json"
	"github.com/ak98neon/KafkaUI/src/kafka"
	"log"
	"net/http"
)

func ConfigHandlerRest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsConfigured {
			errorMsg := "Service isn't configured"
			log.Println(errorMsg)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("400 " + errorMsg))
			return
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}

func GetKafkaInfo(w http.ResponseWriter, _ *http.Request) {
	info := kafka.GetKafkaInfo()
	jsonInfo, err := json.Marshal(info)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonInfo)
}
