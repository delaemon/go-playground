package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("pingHandler", func(t *testing.T) {
		t.Parallel()

		s := httptest.NewServer(http.HandlerFunc(pingHandler()))
		defer s.Close()

		res, err := http.Get(s.URL)
		if err != nil {
			t.Fatal(err)
		}
		if res.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Fatal("error content type: %s", res.Header.Get("Content-Type"))
		}
		if res.StatusCode != 200 {
			t.Fatal("error status code: %d", res.StatusCode)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "pong" {
			t.Fatal("error body: %s", string(body))
		}
	})

	t.Run("echoHandler", func(t *testing.T) {
		candidates := []struct {
			query    string
			expected string
		}{
			{"", ""},
			{"foo=bar", ""},
			{"msg=foo", "foo"},
		}
		for _, c := range candidates {
			c := c
			t.Run(c.query, func(t *testing.T) {
				t.Parallel()

				s := httptest.NewServer(http.HandlerFunc(echoHandler()))
				defer s.Close()

				res, err := http.Get(fmt.Sprintf("%v?%v", s.URL, c.query))
				if err != nil {
					t.Fatal(err)
				}
				if res.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
					t.Fatal("error content type: %s", res.Header.Get("Content-Type"))
				}
				if res.StatusCode != 200 {
					t.Fatal("error status code: %d", res.StatusCode)
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				if string(body) != c.expected {
					t.Fatal("error body: %s", string(body))
				}
			})
		}
	})
}

func TestHandlerWithRecorder(t *testing.T) {
	t.Run("echoHandler", func(t *testing.T) {
		candidates := []struct {
			url      string
			expected string
		}{
			{"http://example.com/?", ""},
			{"http://example.com/?foo=bar", ""},
			{"http://example.com/?msg=foo", "foo"},
		}

		for _, c := range candidates {
			c := c
			t.Run(c.url, func(t *testing.T) {
				t.Parallel()

				res := httptest.NewRecorder()
				req, err := http.NewRequest(http.MethodGet, c.url, nil)
				if err != nil {
					t.Fatal(err)
				}

				handler := echoHandler()
				handler(res, req)

				if res.HeaderMap.Get("Content-Type") != "text/plain; charset=utf-8" {
					t.Fatal("error content type: %s", res.HeaderMap.Get("Content-Type"))
				}
				if res.Code != 200 {
					t.Fatal("error status code: %d", res.Code)
				}

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				if string(body) != c.expected {
					t.Fatal("error body: %s", string(body))
				}
			})
		}
	})
}
