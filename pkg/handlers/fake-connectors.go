// +build fake

package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/luigizuccarelli/golang-gitwebhook-service/pkg/connectors"
	"github.com/microlib/simple"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Inject (force) readAll test error")
}

type FakeConnectors struct {
	Logger *simple.Logger
	Http   *http.Client
	Name   string
	Force  string
}

// Do - used for testing
func (c *FakeConnectors) Do(req *http.Request) (*http.Response, error) {
	if c.Force == "true" {
		return nil, errors.New("forced http error")
	}
	return c.Http.Do(req)
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// logger wrapper
func (r *FakeConnectors) Error(msg string, val ...interface{}) {
	r.Logger.Error(fmt.Sprintf(msg, val...))
}

func (r *FakeConnectors) Info(msg string, val ...interface{}) {
	r.Logger.Info(fmt.Sprintf(msg, val...))
}

func (r *FakeConnectors) Debug(msg string, val ...interface{}) {
	r.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (r *FakeConnectors) Trace(msg string, val ...interface{}) {
	r.Logger.Trace(fmt.Sprintf(msg, val...))
}

// NewTestConnectors - inject our test connectors
func NewTestConnectors(filename string, code int, force string, logger *simple.Logger) connectors.Clients {

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error(fmt.Sprintf("file data %v\n", err))
		panic(err)
	}
	httpclient := NewHttpTestClient(func(r *http.Request) *http.Response {
		logger.Trace(fmt.Sprintf("Request Object %v", r))
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(data))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conn := &FakeConnectors{Http: httpclient, Logger: logger, Force: force, Name: "FakeConnectors"}
	return conn
}
