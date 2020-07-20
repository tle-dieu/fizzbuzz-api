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

	// "github.com/golang/protobuf/proto"
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

func FizzbuzzCheckEncoding(acceptEncoding string) (string, error) {
	contentTypeResponse := [...]string{"application/json", "application/protobuf", "application/xml"}
	arrAcceptEncoding := strings.Split(acceptEncoding, ",")

	if len(arrAcceptEncoding) == 0 || strings.TrimSpace(arrAcceptEncoding[0]) == "" {
		return contentTypeResponse[0], nil
	}
	for _, v := range arrAcceptEncoding {
		for _, contentType := range contentTypeResponse {
			if strings.TrimSpace(v) == contentType {
				return contentType, nil
			}
		}
	}
	return "", fmt.Errorf("Bad Accept-encoding, can be %v", contentTypeResponse)
}

//add protobuf implementation
//http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
func FizzbuzzGetData(req *http.Request) (Data, error) {
	var d Data
	var err error

	if req.Body == nil || req.Body == http.NoBody {
		return d, errors.New("Please send a request body")
	}
	contentType := req.Header.Get("Content-Type")
	if contentType == "application/json" {
		err = json.NewDecoder(req.Body).Decode(&d)
	} else if contentType == "application/protobuf" {
		fmt.Println("decode protobuf")
		// err = proto.Unmarshal(req.Body, &d)
	} else {
		fmt.Println("NO")
	}
	if err != nil {
		return d, err
	}
	return d, nil
}

func FizzbuzzHandle(w http.ResponseWriter, req *http.Request) (int, error) {
	d, err := FizzbuzzGetData(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return http.StatusBadRequest, err
	}

	log.Println("Request Data: ", d)
	if err := FizzbuzzCheckData(d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return http.StatusBadRequest, err
	}

	contentType, err := FizzbuzzCheckEncoding(req.Header.Get("Accept-encoding"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return http.StatusNotAcceptable, err
	}
	w.Header().Set("Content-type", contentType)
	fmt.Fprintln(w, FizzbuzzAlgo(d))
	return http.StatusOK, err
}

func FizzbuzzLog(w http.ResponseWriter, req *http.Request) {
	log.Println("Local Timestamp: ", time.Now())
	log.Println("Request Method: ", req.Method)
	log.Println("Request Url: ", req.URL)
	log.Println("Request Header: ", req.Header)
	statusCode, err := FizzbuzzHandle(w, req)
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
