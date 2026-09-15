package model

// InyeccionRequest represents the incoming request to inject a loan application.
type InyeccionRequest struct {
	DatosCredito  DatosCredito  `json:"datos_credito"`
	Participantes []Participante `json:"participantes"`
	Propiedades   []Propiedad   `json:"propiedades,omitempty"`
}

// InyeccionResponse represents the API response.
type InyeccionResponse struct {
	StatusCode      *int    `json:"StatusCode"`
	Mensaje         string  `json:"Mensaje"`
	NumeroSolicitud *int    `json:"NumeroSolicitud"`
}

// DatosCredito holds the credit details.
type DatosCredito struct {
	NumeroSolicitud    *int     `json:"NumeroSolicitud"`
	AntiguedadVivienda *int     `json:"AntiguedadVivienda"`
	EjecutivoComercial *string  `json:"EjecutivoComercial"`
	Producto           *int     `json:"Producto"`
	Objetivo           *int     `json:"Objetivo"`
	Destino            *int     `json:"Destino"`
	MontoAprobado      *float64 `json:"MontoAprobado"`
	ValorPropiedad     *float64 `json:"ValorPropiedad"`
	FechaAprobacion    *string  `json:"FechaAprobacion"`
	ValorContado       *float64 `json:"ValorContado"`
	Plazo1             *int     `json:"Plazo1"`
	Plazo2             *int     `json:"Plazo2,omitempty"`
	MesesGracia        *int     `json:"MesesGracia"`
	Tasa1              *float64 `json:"Tasa1"`
	Spread1            *float64 `json:"Spread1"`
}

// Participante holds participant details. 
// Adding basic fields, more can be added if present in Java.
type Participante struct {
	Rut      *string `json:"Rut,omitempty"`
	Nombre   *string `json:"Nombre,omitempty"`
	Apellido *string `json:"Apellido,omitempty"`
	Rol      *string `json:"Rol,omitempty"`
}

// Propiedad holds property details.
type Propiedad struct {
	Direccion *string `json:"Direccion,omitempty"`
	Comuna    *string `json:"Comuna,omitempty"`
	Region    *string `json:"Region,omitempty"`
}

// TokenRequest represents the request to the authentication service.
type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse represents the authentication response.
type TokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh,omitempty"`
}
