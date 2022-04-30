package rest

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"runtime/debug"
	"time"

	"go_rest/internal/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		logger.Debugf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Error, Упал сервер", http.StatusInternalServerError)
				log.Println(err)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func validateRequestType(w http.ResponseWriter, req *http.Request, request_type string) error {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if mediatype != request_type {
		http.Error(w, fmt.Sprintf("expect %v Content-Type", request_type), http.StatusUnsupportedMediaType)
		return err
	}
	return nil
}

func ValidateRequestJsonType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := validateRequestType(w, req, "application/json")
		if err != nil {
			return
		}
		next.ServeHTTP(w, req)
	})
}
