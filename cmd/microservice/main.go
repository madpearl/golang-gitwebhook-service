//go:build real
// +build real

package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/luigizuccarelli/golang-gitwebhook-service/pkg/connectors"
	"github.com/luigizuccarelli/golang-gitwebhook-service/pkg/handlers"
	"github.com/luigizuccarelli/golang-gitwebhook-service/pkg/validator"
	"github.com/microlib/simple"
)

func startHttpServer(con connectors.Clients) (*http.Server, error) {
	srv := &http.Server{Addr: ":9000"}
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/service", func(w http.ResponseWriter, r *http.Request) {
		handlers.WebhookHandler(w, r, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/isalive", func(w http.ResponseWriter, r *http.Request) {
		handlers.IsAlive(w, r, con)
	}).Methods("GET", "OPTIONS")

	http.Handle("/", r)

	if err := srv.ListenAndServe(); err != nil {
		con.Error("Httpserver: ListenAndServe() error: %v", err)
		return srv, err
	}

	return srv, nil
}

func main() {
	var logger *simple.Logger

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}

	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}

	conn := connectors.NewClientConnectors(logger)
	srv, err := startHttpServer(conn)
	if err != nil {
		os.Exit(-1)
	}
	logger.Info("Starting server on port " + srv.Addr)
}
