# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instalar certificados CA (necesarios para llamadas HTTPS externas)
RUN apk --no-cache add ca-certificates

# Descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la API quitando símbolos de depuración (-w -s) para hacer el binario extra ligero
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/api ./cmd/api

# Stage 2: Producción (Scratch)
FROM scratch

# Copiar certificados desde el builder para poder hacer peticiones a FinnFlow
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar el ejecutable ultra-ligero
COPY --from=builder /app/api /api

EXPOSE 8085

# Ejecutar el binario
ENTRYPOINT ["/api"]
