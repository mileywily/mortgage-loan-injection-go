package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	myendpoint "cl.bancofalabella.mortgage/injection/internal/endpoint"
	"cl.bancofalabella.mortgage/injection/internal/handler"
	"cl.bancofalabella.mortgage/injection/internal/model"
)

func mockEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	status := 200
	return myendpoint.InjectResponse{
		Response: &model.InyeccionResponse{
			StatusCode: &status,
			Mensaje:    "Success",
		},
	}, nil
}

func TestMakeHTTPHandler(t *testing.T) {
	h := handler.MakeHTTPHandler(mockEndpoint)
	req := httptest.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/injections", bytes.NewBufferString(`{"test":"1"}`))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}
}
