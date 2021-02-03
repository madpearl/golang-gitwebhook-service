// +build fake

package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microlib/simple"
)

// Makes use of the fake-connectors.go (in this package) for testing via the +build fake directive

func TestHandlers(t *testing.T) {

	var logger = &simple.Logger{Level: "info"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("SimpleHandler : should pass (post) merge", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/merge.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SimpleHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SimpleHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SimpleHandler : should pass (post) uat", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/uat-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SimpleHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SimpleHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SimpleHandler : should pass (post) prod", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/prod-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SimpleHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SimpleHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SimpleHandler : should fail (body content error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", errReader(0))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SimpleHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SimpleHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SimpleHandler : should fail (json to golang struct error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte("{ bad json")))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SimpleHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SimpleHandler ", rr.Code, STATUS))
		}
	})

	t.Run("SimpleHandler : should fail (force http error)", func(t *testing.T) {
		var STATUS int = 500

		requestPayload, _ := ioutil.ReadFile("../../tests/prod-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "true", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SimpleHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SimpleHandler ", rr.Code, STATUS))
		}
	})

}
