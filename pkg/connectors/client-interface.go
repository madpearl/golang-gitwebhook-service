package connectors

import "net/http"

// Clients interface - the NewClientConnectors function will implement this interface
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Do(req *http.Request) (*http.Response, error)
}
