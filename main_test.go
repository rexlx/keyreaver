package main

import (
	"bytes"
	"encoding/json"
	"log"
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
			log.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(out)))
		rr := httptest.NewRecorder()
		handler(rr, req)

		if got := rr.Body.String(); got != tc.key {
			t.Errorf("%v : got %v, wanted %v", tc.name, got, tc.key)
		}
	}
}
