// +build real

package connectors

import (
	"crypto/tls"
	"net/http"

	"github.com/microlib/simple"
)

// The premise here is to use this as a reciever in the relevant functions
// this allows us then to mock/fake connections and calls
type Connectors struct {
	Logger *simple.Logger
	Http   *http.Client
	Name   string
}

// NewClientConnectors : function that initialises connections to DB's, caches' queues etc
// Seperating this functionality here allows us to inject a fake or mock connection object for testing
func NewClientConnectors(logger *simple.Logger) Clients {

	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	return &Connectors{Http: httpClient, Logger: logger, Name: "RealConnectors"}
}
