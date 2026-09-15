package service

import (
	"fmt"

	"strings"

	"cl.bancofalabella.mortgage/injection/internal/model"
)

// Custom error types for specific HTTP mappings
type ValidationError struct {
	Code    string
	Message string
	Status  int
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(field string, isFormat bool) *ValidationError {
	if isFormat {
		return &ValidationError{
			Code:    "invalid_format",
			Message: fmt.Sprintf("%s with invalid format", field),
			Status:  402, // PAYMENT_REQUIRED
		}
	}
	return &ValidationError{
		Code:    "invalid_value",
		Message: fmt.Sprintf("%s with invalid value", field),
		Status:  406, // NOT_ACCEPTABLE
	}
}

func NewMissingFieldError(field string) *ValidationError {
	return &ValidationError{
		Code:    "missing_field",
		Message: fmt.Sprintf("Missing %s field", field),
		Status:  400, // BAD_REQUEST
	}
}

// ValidateBusinessRules replicates the logic of InjectionService.validateBusinessRules
func ValidateBusinessRules(req *model.InyeccionRequest) error {
	if req.Participantes == nil || len(req.Participantes) == 0 {
		return NewMissingFieldError("participantes")
	}

	dc := req.DatosCredito
	
	if dc.MontoAprobado != nil && dc.ValorPropiedad != nil && *dc.MontoAprobado > *dc.ValorPropiedad {
		return NewValidationError("MontoAprobado", false)
	}

	if dc.Plazo1 != nil {
		if *dc.Plazo1 < 5 || *dc.Plazo1 > 30 || *dc.Plazo1%5 != 0 {
			return NewValidationError("Plazo1", false)
		}
	}

	if dc.Tasa1 != nil {
		if *dc.Tasa1 < 0.1 || *dc.Tasa1 > 50.0 {
			return NewValidationError("Tasa1", false)
		}
	}
	
	if dc.Spread1 != nil {
		if *dc.Spread1 < 0.0 {
			return NewValidationError("Spread1", false)
		}
	}

	return nil
}

// ExtractErrorMessage replicates InjectionService.extractErrorMessage
func ExtractErrorMessage(err error) string {
	if err == nil {
		return "Error interno del sistema FinnFlow."
	}
	
	msg := err.Error()

	if strings.Contains(msg, "token_not_valid") || strings.Contains(msg, "Token is invalid") {
		return "Token de autenticación no válido para el sistema FinnFlow."
	}
	if strings.Contains(msg, "Authentication failed") {
		return "Error de autenticación: credenciales inválidas para el sistema FinnFlow."
	}
	if strings.Contains(msg, "401") {
		return "Error de autenticación con el sistema FinnFlow."
	}
	if strings.Contains(msg, "400") {
		return "Datos de solicitud no válidos para el sistema FinnFlow."
	}
	if strings.Contains(msg, "403") {
		return "Acceso denegado al sistema FinnFlow."
	}
	if strings.Contains(msg, "404") {
		return "Servicio no encontrado en el sistema FinnFlow."
	}
	if strings.Contains(msg, "timeout") || strings.Contains(msg, "connect") || strings.Contains(msg, "context deadline exceeded") {
		return "Timeout en la comunicación con FinnFlow."
	}
	if strings.Contains(msg, "certificate") || strings.Contains(msg, "SSL") || strings.Contains(msg, "PKIX") {
		return "Error de conectividad SSL con FinnFlow."
	}

	return "Error interno del sistema FinnFlow."
}
