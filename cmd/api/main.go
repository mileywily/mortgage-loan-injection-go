package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"

	"cl.bancofalabella.mortgage/injection/internal/activities"
	"cl.bancofalabella.mortgage/injection/internal/endpoint"
	"cl.bancofalabella.mortgage/injection/internal/handler"
	"cl.bancofalabella.mortgage/injection/internal/workflows"
	"cl.bancofalabella.mortgage/injection/pkg/logger"
)

func main() {
	logger.InitLogger()

	finnflowURL := os.Getenv("FINNFLOW_URL")
	if finnflowURL == "" {
		finnflowURL = "http://localhost:8081" // fallback for local mock
	}
	username := os.Getenv("FINNFLOW_KEY")
	password := os.Getenv("FINNFLOW_SECRET")

	// Dependencies
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if os.Getenv("INSECURE_SKIP_VERIFY") == "true" {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		slog.Warn("WARNING: INSECURE_SKIP_VERIFY is set to true. SSL certificate validation is disabled.")
	}
	httpClient := &http.Client{Transport: customTransport}
	finnFlowClient := activities.NewFinnFlowClient(httpClient, finnflowURL)
	workflow := workflows.NewInjectionWorkflow(finnFlowClient, username, password)
	ep := endpoint.MakeInjectEndpoint(workflow)
	
	httpHandler := handler.MakeHTTPHandler(ep)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Starting server on port " + port)
	if err := http.ListenAndServe(":"+port, httpHandler); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
