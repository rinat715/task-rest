package rest

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"runtime/debug"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func ValidateRequestType(w http.ResponseWriter, req *http.Request, request_type string) error {
	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if mediatype != request_type {
		http.Error(w, fmt.Sprint("expect %s Content-Type", request_type), http.StatusUnsupportedMediaType)
		return err
	}
	return nil
}
