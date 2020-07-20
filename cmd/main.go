package main

import (
	"log"
	"net/http"
	"time"

	// "github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/tle-dieu/fizzbuzz-api/internal/handler"
	"github.com/tle-dieu/fizzbuzz-api/internal/logger"
)

func main() {
	router := mux.NewRouter()
	serv := &http.Server{
		Handler:      logger.WrapHandlerWithLogging(router),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	router.HandleFunc("/fizzbuzz", handler.Fizzbuzz).Methods(http.MethodPost)
	log.Println("Listening on :8080")
	log.Fatal(serv.ListenAndServe())
}
