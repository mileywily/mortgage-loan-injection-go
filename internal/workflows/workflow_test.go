package workflows_test

import (
	"context"
	"testing"

	"cl.bancofalabella.mortgage/injection/internal/model"
	"cl.bancofalabella.mortgage/injection/internal/workflows"
)

type mockClient struct {
	tokenErr error
	injErr   error
	injResp  *model.InyeccionResponse
}

func (m *mockClient) GetAccessToken(ctx context.Context, u, p string) (string, error) {
	return "token123", m.tokenErr
}
func (m *mockClient) InjectToSalesforce(ctx context.Context, r, t string) (*model.InyeccionResponse, error) {
	return m.injResp, m.injErr
}

func TestProcessInjection_Success(t *testing.T) {
	statusCode := 200
	mc := &mockClient{
		injResp: &model.InyeccionResponse{StatusCode: &statusCode, Mensaje: "OK"},
	}
	wf := workflows.NewInjectionWorkflow(mc, "u", "p")

	resp, err := wf.ProcessInjection(context.Background(), "req")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", *resp.StatusCode)
	}
}
