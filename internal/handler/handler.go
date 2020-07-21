package handler

import (
	"fmt"
	"net/http"

	"github.com/tle-dieu/fizzbuzz-api/pkg/fizzbuzz"
)

//Fizzbuzz handle a fizzbuzz request (/fizzbuzz)
func Fizzbuzz(w http.ResponseWriter, req *http.Request) {
	d, err := fizzbuzzGetDataRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := fizzbuzz.CheckData(d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType, err := fizzbuzzCheckResponseType(req.Header.Get("Accept-encoding"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	w.Header().Set("Content-type", contentType)
	fmt.Fprintln(w, fizzbuzz.FizzbuzzAlgo(d))
}
