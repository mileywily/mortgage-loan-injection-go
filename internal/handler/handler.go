package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"cl.bancofalabella.mortgage/injection/docs"
	myendpoint "cl.bancofalabella.mortgage/injection/internal/endpoint"
	"cl.bancofalabella.mortgage/injection/internal/model"
	"cl.bancofalabella.mortgage/injection/internal/service"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// MakeHTTPHandler creates the HTTP handler for the injection service.
func MakeHTTPHandler(e endpoint.Endpoint, options ...httptransport.ServerOption) http.Handler {
	m := http.NewServeMux()
	
	// Ensure we handle custom errors effectively
	options = append(options, httptransport.ServerErrorEncoder(encodeError))

	// API Endpoints
	m.Handle("/v1/bfcl/mortgage-loan/injections", httptransport.NewServer(
		e,
		decodeInjectRequest,
		encodeInjectResponse,
		options...,
	))

	// Swagger Endpoints
	m.HandleFunc("/api-docs", func(w http.ResponseWriter, r *http.Request) {
		data, _ := docs.Files.ReadFile("swagger.json")
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	m.HandleFunc("/swagger-ui.html", func(w http.ResponseWriter, r *http.Request) {
		data, _ := docs.Files.ReadFile("swagger-ui.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})

	return m
}

func decodeInjectRequest(_ context.Context, r *http.Request) (interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, service.NewMissingFieldError("body")
	}
	return myendpoint.InjectRequest{RawRequest: string(body)}, nil
}

func encodeInjectResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(myendpoint.InjectResponse)
	
	if resp.Err != nil {
		encodeError(ctx, resp.Err, w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	// If the inner response is of type InyeccionResponse, we should check its StatusCode
	// just like the Java controller did:
	// if (response.getStatusCode() != null && response.getStatusCode() >= 400) { return ResponseEntity.status... }
	if inyResp, ok := resp.Response.(*model.InyeccionResponse); ok {
		if inyResp.StatusCode != nil && *inyResp.StatusCode >= 400 {
			w.WriteHeader(*inyResp.StatusCode)
			return json.NewEncoder(w).Encode(inyResp)
		}
	}
	
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp.Response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	if valErr, ok := err.(*service.ValidationError); ok {
		w.WriteHeader(valErr.Status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    valErr.Code,
			"message": valErr.Message,
		})
		return
	}
	
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    "internal_error",
		"message": "Internal server error",
	})
}
