package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Data struct {
	Str1  string
	Str2  string
	Int1  int
	Int2  int
	Limit int
}

type StringError struct {
	msg string
}

func FizzbuzzAlgo(d Data) string {
	var response bytes.Buffer
	var write bool

	for i := 1; i < d.Limit; i++ {
		write = false
		if i%d.Int1 == 0 {
			response.WriteString(d.Str1)
			write = true
		}
		if i%d.Int2 == 0 {
			response.WriteString(d.Str2)
			write = true
		}
		if !write {
			response.WriteString(strconv.Itoa(i))
		}
		response.WriteString(" ")
	}
	return response.String()
}

func FizzbuzzCheckData(d Data) error {
	if d.Int1 <= 0 {
		return errors.New("int1 should by greater than zero")
	}
	if d.Int2 <= 0 {
		return errors.New("int2 should by greater than zero")
	}
	if d.Str1 == "" {
		return errors.New("str1 is empty")
	}
	if d.Str2 == "" {
		return errors.New("str2 is empty")
	}
	return nil
}

func FizzbuzzCheckEncoding(acceptEncoding string) (error, string) {
	contentTypeResponse := [...]string{"application/json", "application/protobuf", "application/xml"}
	arrAcceptEncoding := strings.Split(acceptEncoding, ",")

	if len(arrAcceptEncoding) == 0 || strings.TrimSpace(arrAcceptEncoding[0]) == "" {
		return nil, contentTypeResponse[0]
	}
	for _, v := range arrAcceptEncoding {
		for _, contentType := range contentTypeResponse {
			if strings.TrimSpace(v) == contentType {
				return nil, contentType
			}
		}
	}
	return errors.New(fmt.Sprintf("Bad Accept-encoding, can be %v", contentTypeResponse)), ""
}

func FizzbuzzGetData(req *http.Request) (error, Data) {
	var d Data

	if req.Body == nil || req.Body == http.NoBody {
		return errors.New("Please send a request body"), d
	}
	if req.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type should be 'application/json'"), d
	}
	if err := json.NewDecoder(req.Body).Decode(&d); err != nil {
		return err, d
	}
	return nil, d
}

func FizzbuzzHandle(w http.ResponseWriter, req *http.Request) {
	err, d := FizzbuzzGetData(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("data request:", d)
	if err := FizzbuzzCheckData(d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err, contentType := FizzbuzzCheckEncoding(req.Header.Get("Accept-encoding"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	w.Header().Set("Content-type", contentType)
	fmt.Fprintln(w, FizzbuzzAlgo(d))
}

//        http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
func main() {
	router := mux.NewRouter()
	serv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	handlers.LoggingHandler(os.Stdout, router)
	router.HandleFunc("/fizz-buzz", FizzbuzzHandle).Methods(http.MethodPost)
	log.Println("Listening on :8080")
	log.Fatal(serv.ListenAndServe())
}
