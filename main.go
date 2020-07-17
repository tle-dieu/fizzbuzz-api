package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		return errors.New("int1 should by greater than zero\n")
	}
	if d.Int2 <= 0 {
		return errors.New("int2 should by greater than zero\n")
	}
	if d.Str1 == "" {
		return errors.New("str1 is empty\n")
	}
	if d.Str2 == "" {
		return errors.New("str2 is empty\n")
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
	return errors.New(fmt.Sprintf("Bad Accept-encoding, can be %v\n", contentTypeResponse)), ""
}

//add protobuf implementation
//http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
func FizzbuzzGetData(req *http.Request) (error, Data) {
	var d Data

	if req.Body == nil || req.Body == http.NoBody {
		return errors.New("Please send a request body\n"), d
	}
	if req.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type should be 'application/json'\n"), d
	}
	if err := json.NewDecoder(req.Body).Decode(&d); err != nil {
		return err, d
	}
	return nil, d
}

func FizzbuzzHandle(w http.ResponseWriter, req *http.Request) (error, int) {
	err, d := FizzbuzzGetData(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err, http.StatusBadRequest
	}

	log.Println("Request Data: ", d)
	if err := FizzbuzzCheckData(d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err, http.StatusBadRequest
	}

	err, contentType := FizzbuzzCheckEncoding(req.Header.Get("Accept-encoding"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return err, http.StatusNotAcceptable
	}
	w.Header().Set("Content-type", contentType)
	fmt.Fprintln(w, FizzbuzzAlgo(d))
	return err, http.StatusOK
}

func FizzbuzzLog(w http.ResponseWriter, req *http.Request) {
	log.Println("Local Timestamp: ", time.Now())
	log.Println("Request Method: ", req.Method)
	log.Println("Request Url: ", req.URL)
	log.Println("Request Header: ", req.Header)
	err, statusCode := FizzbuzzHandle(w, req)
	log.Println("Response Status Code: ", statusCode)
	log.Println("Response Header: ", w.Header())
	if err != nil {
		log.Println("Error: ", err.Error())
	}
}

func main() {
	router := mux.NewRouter()
	serv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	router.HandleFunc("/fizz-buzz", FizzbuzzLog).Methods(http.MethodPost)
	log.Println("Listening on :8080")
	log.Fatal(serv.ListenAndServe())
}
