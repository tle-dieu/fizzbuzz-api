package logger

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

//WrapHandlerWithLogging adds logger to a handler
func WrapHandlerWithLogging(wrappedHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		lrw := &loggingResponseWriter{w, http.StatusOK}
		log.Println("Local Timestamp: ", time.Now())
		log.Printf("Request: [method: %s, URL: %s, header: %s]\n", req.Method, req.URL, req.Header)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		wrappedHandler.ServeHTTP(lrw, req)
		log.Printf("Data Request: %q\n", body)
		log.Printf("Response: [header: %s, status: %d]\n", w.Header(), lrw.statusCode)
	}
}
