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
	if req.Header.Get("Content-type") != "application/json" {
		return errors.New("Content-type should be 'application/json'"), d
	}
	if err := json.NewDecoder(req.Body).Decode(&d); err != nil {
		return err, d
	}
	return nil, d
}

func FizzbuzzHandlePost(w http.ResponseWriter, req *http.Request) {

	err, d := FizzbuzzGetData(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println("data request:", d)
	if err := FizzbuzzCheckData(d); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err, contentType := FizzbuzzCheckEncoding(req.Header.Get("Accept-encoding"))
	if err != nil {
		http.Error(w, err.Error(), 406)
		return
	}
	w.Header().Set("Content-type", contentType)
	fmt.Fprintln(w, FizzbuzzAlgo(d))
}

func main() {
	http.HandleFunc("/fizz-buzz", FizzbuzzHandlePost)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
