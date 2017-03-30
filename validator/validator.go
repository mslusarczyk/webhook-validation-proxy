package validator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/rjz/githubhook"
)

type Validator struct {
	secret    string
	sourceNet *net.IPNet
}

func NewValidator(secret, sourceCIDR string) *Validator {
	_, parseCIDR, err := net.ParseCIDR(sourceCIDR)
	if err != nil {
		log.Fatalf("Could not parse CIDR: %s, err: %s", sourceCIDR, err)
	}

	if len(secret) == 0 {
		log.Println("Secret is not set, validation has little sense it this state.")
	}

	return &Validator{secret: secret, sourceNet: parseCIDR}
}

func (v Validator) Validate(req *http.Request) (bool, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)

	if err != nil {
		return false, errors.New(fmt.Sprintf("Could not get remote IP from req, err: %s", err))
	}

	if !v.sourceNet.Contains(net.ParseIP(ip)) {
		return false, errors.New("Remote IP incorrect")
	}

	hook, err := githubhook.Parse([]byte(v.secret), req)

	if err != nil {
		return false, errors.New(fmt.Sprintf("Webhook umarshalling failed, err: %s", err))
	}

	//necessary because githubhook.Parse reads req.Body
	req.Body = ioutil.NopCloser(bytes.NewBuffer(hook.Payload))

	return true, nil
}
