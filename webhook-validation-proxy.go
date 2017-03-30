package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/mslusarczyk/webhook-validation-proxy/config"
	"github.com/mslusarczyk/webhook-validation-proxy/validator"
)

var (
	proxy            *httputil.ReverseProxy
	webhookValidator *validator.Validator
)

func main() {
	params := config.ParseParams()

	webhookValidator = validator.NewValidator(params.Secret, params.Cidr)

	// proxy
	addr, err := url.Parse(params.Target)
	if err != nil {
		log.Fatalf("Could not parse target url, err: %s", err)
	}
	proxy = httputil.NewSingleHostReverseProxy(addr)

	// main server
	http.Handle(params.Context, handlers.LoggingHandler(os.Stdout, http.HandlerFunc(validateAndProxy)))

	http.Handle("/", http.NotFoundHandler())
	http.ListenAndServe(":"+params.Port, nil)
}

func validateAndProxy(w http.ResponseWriter, req *http.Request) {
	_, err := webhookValidator.Validate(req)

	if err != nil {
		msg := "Validation failed"
		log.Printf("%s, err: %s", msg, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
	}

	proxy.ServeHTTP(w, req)
}
