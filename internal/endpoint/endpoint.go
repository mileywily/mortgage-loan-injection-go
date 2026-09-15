package endpoint

import (
	"context"

	"cl.bancofalabella.mortgage/injection/internal/workflows"
	"github.com/go-kit/kit/endpoint"
)

// InjectRequest wraps the raw string request for Go-Kit
type InjectRequest struct {
	RawRequest string
}

// InjectResponse wraps the response
type InjectResponse struct {
	Response interface{}
	Err      error
}

// Failed implements endpoint.Failer
func (r InjectResponse) Failed() error {
	return r.Err
}

// MakeInjectEndpoint creates the Go-Kit endpoint
func MakeInjectEndpoint(wf workflows.InjectionWorkflow) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(InjectRequest)
		resp, err := wf.ProcessInjection(ctx, req.RawRequest)
		
		// If the workflow returns a specific domain error (e.g. ValidationError),
		// we pass it in the response to let the HTTP transport handle status codes.
		if err != nil {
			return InjectResponse{Err: err}, nil
		}
		
		return InjectResponse{Response: resp, Err: nil}, nil
	}
}
