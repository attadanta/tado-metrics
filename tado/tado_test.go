package tado

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// A fake tado service answering an http request.
type tadoService func(*http.Request) (*http.Response, error)

func (t tadoService) RoundTrip(req *http.Request) (*http.Response, error) {
	return t(req)
}

var token = "letmein"

func TestHomeId(t *testing.T) {
	jsonResponse := `{"homeId": 3}`

	c := &http.Client{
		Transport: tadoService(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(jsonResponse)),
			}, nil
		}),
	};

	want := 3
	got := HomeId(c, token)

	if got != want {
		t.Errorf("Got incorrect home id, want %q, got %q", want, got)
	}
}
