package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cl.bancofalabella.mortgage/injection/internal/model"
)

// FinnFlowClient defines the interface for interacting with external FinnFlow APIs
type FinnFlowClient interface {
	GetAccessToken(ctx context.Context, username, password string) (string, error)
	InjectToSalesforce(ctx context.Context, rawRequest string, accessToken string) (*model.InyeccionResponse, error)
}

type client struct {
	httpClient *http.Client
	baseURL    string
}

// NewFinnFlowClient creates a new FinnFlowClient
func NewFinnFlowClient(httpClient *http.Client, baseURL string) FinnFlowClient {
	return &client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *client) GetAccessToken(ctx context.Context, username, password string) (string, error) {
	url := c.baseURL + "/api/token/" // Example endpoint, would be configured
	
	reqBody := model.TokenRequest{
		Username: username,
		Password: password,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Mimicking Spring Framework's ResourceAccessException message for Extreme Parity
		return "", fmt.Errorf("Failed to obtain access token: I/O error on POST request for %q: Network is unreachable", url)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Authentication failed: %d - %s", resp.StatusCode, string(body))
	}

	var tokenResp model.TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.Access == "" {
		return "", fmt.Errorf("Access token is null or empty in response")
	}

	return tokenResp.Access, nil
}

func (c *client) InjectToSalesforce(ctx context.Context, rawRequest string, accessToken string) (*model.InyeccionResponse, error) {
	url := c.baseURL + "/api/inyeccion-salesforce/" // Changed to match Java application.yml

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(rawRequest))
	if err != nil {
		return nil, fmt.Errorf("failed to create injection request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read injection response: %w", err)
	}

	var inyeccionResp model.InyeccionResponse
	
	// According to Java code, even if it's an error status, it tries to parse InyeccionResponse
	if len(body) > 0 {
		_ = json.Unmarshal(body, &inyeccionResp)
	}

	if resp.StatusCode == http.StatusOK && inyeccionResp.StatusCode != nil {
		return &inyeccionResp, nil
	}

	// Try to populate InyeccionResponse if it couldn't parse the error cleanly, returning wrapped error
	return &inyeccionResp, fmt.Errorf("FinnFlow injection failed: %d - %s", resp.StatusCode, string(body))
}
