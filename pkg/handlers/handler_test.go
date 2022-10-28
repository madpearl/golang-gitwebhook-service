//go:build fake
// +build fake

package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/microlib/simple"
)

// Makes use of the fake-connectors.go (in this package) for testing via the +build fake directive

func TestHandlers(t *testing.T) {

	var logger = &simple.Logger{Level: "info"}

	os.Setenv("PR_OPENED_URL", "loclahost")
	os.Setenv("PR_MERGED_URL", "loclahost")
	os.Setenv("PRERELEASED_URL", "localhost")
	os.Setenv("RELEASED_URL", "localhost")

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/isalive", nil)
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			IsAlive(w, r, conn)
		})
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

	t.Run("WebhookHandler : should pass (post) pr", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/git-payload.txt")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebhookHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "WebhookHandler ", rr.Code, STATUS))
		}
	})

	t.Run("WebhookHandler : should pass (post) publish", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/git-payload-published.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebhookHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "WebhookHandler ", rr.Code, STATUS))
		}
	})

	t.Run("WebhookHandler : should pass (post) uat", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/uat-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebhookHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "WebhookHandler ", rr.Code, STATUS))
		}
	})

	t.Run("WebhookHandler : should pass (post) prod", func(t *testing.T) {
		var STATUS int = 200

		requestPayload, _ := ioutil.ReadFile("../../tests/prod-release.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebhookHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "WebhookHandler ", rr.Code, STATUS))
		}
	})

	t.Run("WebhookHandler : should fail (body content error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", errReader(0))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebhookHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "WebhookHandler ", rr.Code, STATUS))
		}
	})

	t.Run("WebhookHandler : should fail (json to golang struct error)", func(t *testing.T) {
		var STATUS int = 500

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte("{ bad json")))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "none", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebhookHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "WebhookHandler ", rr.Code, STATUS))
		}
	})

	t.Run("WebhookHandler : should fail (force http error)", func(t *testing.T) {
		var STATUS int = 500

		requestPayload, _ := ioutil.ReadFile("../../tests/git-payload-pr-merged.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/service", bytes.NewBuffer([]byte(requestPayload)))
		conn := NewTestConnectors("../../tests/response.json", STATUS, "true", logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebhookHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "WebhookHandler ", rr.Code, STATUS))
		}
	})

}
