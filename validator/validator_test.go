package validator_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
	"testing"

	"github.com/mslusarczyk/webhook-validation-proxy/validator"
)

const (
	testSecret = "shhhhhhh"
	testCIDR   = "192.30.152.0/22"
)

func TestNoIPPost(t *testing.T) {
	request, _ := http.NewRequest("GET", "/path", nil)
	validate, error := validate(request)
	assertInvalid("Could not get remote IP", validate, error, t)
}

func TestNonPost(t *testing.T) {
	request, _ := http.NewRequest("GET", "/path", nil)
	request.RemoteAddr = "192.30.153.0:8080"
	validate, error := validate(request)
	assertInvalid("Unknown method", validate, error, t)
}

func TestMissingSignature(t *testing.T) {
	request, _ := http.NewRequest("POST", "/path", nil)
	request.RemoteAddr = "192.30.153.0:8080"
	validate, error := validate(request)
	assertInvalid("umarshalling failed", validate, error, t)
}

func TestMissingEvent(t *testing.T) {
	request, _ := http.NewRequest("POST", "/path", nil)
	request.RemoteAddr = "192.30.153.0:8080"
	request.Header.Add("x-hub-signature", "some signature")
	validate, error := validate(request)
	assertInvalid("No event", validate, error, t)
}

func TestMissingEventId(t *testing.T) {
	request, _ := http.NewRequest("POST", "/path", nil)
	request.RemoteAddr = "192.30.153.0:8080"
	request.Header.Add("x-hub-signature", "some signature")
	request.Header.Add("x-github-event", "some event")
	validate, error := validate(request)
	assertInvalid("No event Id", validate, error, t)
}

func TestInvalidSignature(t *testing.T) {
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("..."))
	r.RemoteAddr = "192.30.153.0:8080"
	r.Header.Add("x-hub-signature", "some signature")
	r.Header.Add("x-github-event", "some event")
	r.Header.Add("x-github-delivery", "some id")

	validate, error := validate(r)
	assertInvalid("Invalid signature", validate, error, t)
}

func TestValidSignature(t *testing.T) {

	body := "{}"

	r, _ := http.NewRequest("POST", "/path", strings.NewReader(body))
	r.RemoteAddr = "192.30.153.0:80"
	r.Header.Add("x-hub-signature", signature(body))
	r.Header.Add("x-github-event", "some event")
	r.Header.Add("x-github-delivery", "some id")

	validate, error := validate(r)
	if error != nil || !validate {
		t.Fatalf("Validation failed, err: %s", error)
	}
}

func TestMultipleValidation(t *testing.T) {
	body := "{}"

	r, _ := http.NewRequest("POST", "/path", strings.NewReader(body))
	r.RemoteAddr = "192.30.153.0:80"
	r.Header.Add("x-hub-signature", signature(body))
	r.Header.Add("x-github-event", "some event")
	r.Header.Add("x-github-delivery", "some id")

	validate(r)
	validate(r)
	validate, error := validate(r)

	if error != nil || !validate {
		t.Fatalf("Validation failed, err: %s", error)
	}
}

func signature(body string) string {
	result := make([]byte, 40)
	computed := hmac.New(sha1.New, []byte(testSecret))
	computed.Write([]byte(body))
	hex.Encode(result, computed.Sum(nil))
	return "sha1=" + string(result)
}

func validate(r *http.Request) (bool, error) {
	underTest := validator.NewValidator(testSecret, testCIDR)
	return underTest.Validate(r)
}

func assertInvalid(msg string, value bool, err error, t *testing.T) {
	if err != nil && strings.Contains(err.Error(), msg) {
		return
	}
	if !value {
		return
	}
	t.Fatalf("Validation failed, expected: [%s] but received [%s]", msg, err)
}
