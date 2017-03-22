package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/namsral/flag"
	"github.com/rjz/githubhook"
)

var (
	context string
	port    string
	target  string
	secret  string
	proxy   *httputil.ReverseProxy
)

func main() {
	// params
	flag.String(flag.DefaultConfigFlagname, "", "Path to config file")
	flag.StringVar(&context, "context", "/github-webhook/", "Context path for proxy for webhooks handling")
	flag.StringVar(&port, "port", "8888", "Port for proxy to listen on")
	flag.StringVar(&target, "target", "http://localhost:8080", "Target address with port")
	flag.StringVar(&secret, "secret", "", "Secret assosiated with GH webhook")
	flag.Parse()

	if len(secret) == 0 {
		log.Println("Secret is not set, validation has little sense it this state.")
	}

	// proxy
	addr, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Could not parse target url, err: %s", err)
	}
	proxy = httputil.NewSingleHostReverseProxy(addr)

	// main server
	http.Handle(context, handlers.LoggingHandler(os.Stdout, http.HandlerFunc(validateAndProxy)))

	http.Handle("/", http.NotFoundHandler())
	http.ListenAndServe(":"+port, nil)
}

func validateAndProxy(w http.ResponseWriter, req *http.Request) {
	hook, err := githubhook.Parse([]byte(secret), req)

	if err != nil {
		msg := "Hook invalid"
		log.Printf("%s, err: %s", msg, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	log.Printf("Handling event: %s", hook.Event)

	//necessary because githubhook.Parse reads req.Body
	req.Body = ioutil.NopCloser(bytes.NewBuffer(hook.Payload))

	proxy.ServeHTTP(w, req)
}
