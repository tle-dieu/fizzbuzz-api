package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tle-dieu/fizzbuzz-api/pkg/fizzbuzz"
	"net/http"
	"strings"
	// "github.com/golang/protobuf/proto"
)

func fizzbuzzCheckResponseType(acceptEncoding string) (string, error) {
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
func FizzbuzzGetDataRequest(req *http.Request) (fizzbuzz.Data, error) {
	var d fizzbuzz.Data
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
		err = errors.New("Content-Type must be application/json or application/protobuf")
	}
	if err != nil {
		return d, err
	}
	return d, nil
}
