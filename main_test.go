package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	os.Setenv("CHALLENGE", "deny all")
	os.Setenv("STORAGE", "nullferatu")
	os.Setenv("FNAME", "tests.txt")
	tests := []struct {
		name   string
		phrase string
		key    string
	}{
		{
			name:   "authorized",
			phrase: "gg7gf7f7vgg",
			key:    "key reaved",
		},
		{
			name:   "unauthorized",
			phrase: "hehe",
			key:    "Service Unavailable",
		},
	}
	for _, tc := range tests {
		var pl payload
		if tc.name == "authorized" {
			os.Setenv("CHALLENGE", tc.phrase)
		}
		pl.TestPhrase = tc.phrase
		pl.Key = "no key here just a test"
		out, err := json.Marshal(pl)
		if err != nil {
			t.Log(err)
			os.Exit(1)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(out)))
		rr := httptest.NewRecorder()
		handler(rr, req)

		if got := rr.Body.String(); got != tc.key {
			t.Errorf("%v : got %v, wanted %v", tc.name, got, tc.key)
		}
	}
}

func Test_readJSON(t *testing.T) {
	var decodedJson struct {
		Foo string `json:"foo"`
	}
	// create sample json
	jason := map[string]interface{}{
		"foo": "bar",
	}
	body, _ := json.Marshal(jason)
	// create req
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	// need  a test recorder
	rr := httptest.NewRecorder()
	defer req.Body.Close()

	err := readJSON(rr, req, &decodedJson)
	if err != nil {
		t.Error("coudlnt decode json", err)
	}

}
