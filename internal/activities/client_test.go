package activities_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cl.bancofalabella.mortgage/injection/internal/activities"
)

func TestFinnFlowClient_GetAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access": "token123"}`))
	}))
	defer ts.Close()

	client := activities.NewFinnFlowClient(ts.Client(), ts.URL)
	token, err := client.GetAccessToken(context.Background(), "user", "pass")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "token123" {
		t.Errorf("expected token123, got %s", token)
	}
}

func TestFinnFlowClient_InjectToSalesforce(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode": 200, "Mensaje": "Success"}`))
	}))
	defer ts.Close()

	client := activities.NewFinnFlowClient(ts.Client(), ts.URL)
	resp, err := client.InjectToSalesforce(context.Background(), `{"payload": "test"}`, "token123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.StatusCode == nil || *resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %v", resp)
	}
}
