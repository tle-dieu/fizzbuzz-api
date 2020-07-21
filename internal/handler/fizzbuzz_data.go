package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tle-dieu/fizzbuzz-api/internal/protobuf"
	"github.com/tle-dieu/fizzbuzz-api/pkg/fizzbuzz"
	"google.golang.org/protobuf/proto"
)

func fizzbuzzCheckResponseType(acceptEncoding string) (string, error) {
	contentTypeResponse := [...]string{"application/json", "application/protobuf", "application/xml"}
	arrAcceptEncoding := strings.Split(acceptEncoding, ",")

	fmt.Println(acceptEncoding)
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

func fizzbuzzProtoToJson(protoReq *protobuf.Data, d *fizzbuzz.Data) {
	d.Str1 = protoReq.GetStr1()
	d.Str2 = protoReq.GetStr2()
	d.Int1 = int(protoReq.GetInt1())
	d.Int2 = int(protoReq.GetInt2())
	d.Limit = int(protoReq.GetLimit())
	fmt.Println(protoReq.GetStr1(), protoReq.GetStr2(), int(protoReq.GetInt1()), int(protoReq.GetInt2()), int(protoReq.GetLimit()))
	fmt.Println(d)
}

func fizzbuzzGetDataRequest(req *http.Request) (fizzbuzz.Data, error) {
	var d fizzbuzz.Data
	var err error

	if req.Body == nil || req.Body == http.NoBody {
		return d, errors.New("Please send a request body")
	}
	contentType := req.Header.Get("Content-Type")
	if contentType == "application/json" {
		err = json.NewDecoder(req.Body).Decode(&d)
	} else if contentType == "application/protobuf" {
		protoReq := &protobuf.Data{}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			err = errors.New("Unable to read message from request : " + err.Error())
		} else {
			err = proto.Unmarshal(body, protoReq)
			fizzbuzzProtoToJson(protoReq, &d)
		}
	} else {
		err = errors.New("Content-Type must be application/json or application/protobuf")
	}
	if err != nil {
		return d, err
	}
	return d, nil
}
