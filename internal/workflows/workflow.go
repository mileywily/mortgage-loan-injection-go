package workflows

import (
	"context"

	"cl.bancofalabella.mortgage/injection/internal/activities"
	"cl.bancofalabella.mortgage/injection/internal/model"
)

// InjectionWorkflow defines the workflow interface
type InjectionWorkflow interface {
	ProcessInjection(ctx context.Context, rawRequest string) (*model.InyeccionResponse, error)
}

type injectionWorkflow struct {
	client activities.FinnFlowClient
	// Other dependencies (like configurations) can go here
	username string
	password string
}

func NewInjectionWorkflow(client activities.FinnFlowClient, username, password string) InjectionWorkflow {
	return &injectionWorkflow{
		client:   client,
		username: username,
		password: password,
	}
}

func (w *injectionWorkflow) ProcessInjection(ctx context.Context, rawRequest string) (*model.InyeccionResponse, error) {
	// Replicating Java's technical debt: Java's InjectionController never actually calls validateBusinessRules
	// and simply forwards the raw request directly to FinnFlow. We bypass validation to maintain Drop-In Replacement.


	// Step 1: Get access token
	accessToken, err := w.client.GetAccessToken(ctx, w.username, w.password)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "" {
			errorMessage = "Error interno del sistema FinnFlow."
		}
		status500 := 500
		return &model.InyeccionResponse{
			StatusCode: &status500,
			Mensaje:    errorMessage,
		}, nil
	}

	// Step 2: Inject to Salesforce
	resp, err := w.client.InjectToSalesforce(ctx, rawRequest, accessToken)
	if err != nil {
		if resp != nil && resp.StatusCode != nil {
			// Forward the parsed error response if available from FinnFlow
			return resp, nil
		}
		errorMessage := err.Error()
		if errorMessage == "" {
			errorMessage = "Error interno del sistema FinnFlow."
		}
		status500 := 500
		return &model.InyeccionResponse{
			StatusCode: &status500,
			Mensaje:    errorMessage,
		}, nil
	}

	return resp, nil
}
