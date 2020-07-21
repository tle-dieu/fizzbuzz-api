package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tle-dieu/fizzbuzz-api/internal/protobuf"
	"google.golang.org/protobuf/proto"
)

func main() {
	d := &protobuf.Data{
		Str1:  "fizz",
		Str2:  "buzz",
		Int1:  3,
		Int2:  6,
		Limit: 20,
	}
	requestBody, err := proto.Marshal(d)
	if err != nil {
		log.Fatalln(err)
	}
	client := new(http.Client)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/fizzbuzz", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Accept-Encoding", "application/json")
	req.Header.Add("Content-Type", "application/protobuf")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
}
