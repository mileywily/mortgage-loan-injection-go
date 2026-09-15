# Mortgage Loan Injection (Go Version)

Este repositorio contiene la migración a Go del microservicio de inyección de créditos hipotecarios hacia FinnFlow.

## Pruebas Locales (Sin necesidad de Docker)

Para los equipos de QA o de Negocio que requieran validar el flujo de la aplicación utilizando un simulador local (Mock Backend / Dummy) sin depender de contenedores Docker ni de entornos Java, existen dos alternativas de ejecución:

### Alternativa 1: Ejecución mediante Script (Para equipos con Go instalado)

Si los evaluadores cuentan con el compilador de Go en sus máquinas, pueden levantar todo el ecosistema (API + Mock) con un solo comando.

**En Windows (PowerShell):**
```powershell
.\scripts\start_dev.ps1
```

**En Mac / Linux (Terminal):**
Pueden ejecutar los servicios abriendo dos pestañas de terminal:
*Terminal 1 (Levanta el Dummy FinnFlow):*
```bash
go run ./cmd/mock_backend
```
*Terminal 2 (Levanta la API Principal):*
```bash
export PORT=8085
export FINNFLOW_URL="http://localhost:8081"
go run ./cmd/api
```

---

### Alternativa 2: Ejecución mediante Binarios Nativos (Para equipos bloqueados corporativamente)

Si los evaluadores (ej. QA) no tienen permisos para instalar Go ni Docker, puedes pre-compilar el código y enviarles **los archivos ejecutables nativos**. Al ser binarios estáticos, correrán instantáneamente con doble clic.

Go permite *Cross-Compilación*, lo que significa que desde tu máquina puedes generar los binarios para el sistema operativo de ellos.

**Si el evaluador usa Windows:**
Ejecuta esto en tu máquina para generar los `.exe`:
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o api.exe ./cmd/api
go build -ldflags="-s -w" -o mock_backend.exe ./cmd/mock_backend
```
*Instrucciones para el QA:* Solo debe darle doble clic a `mock_backend.exe` y luego a `api.exe`.

**Si el evaluador usa Mac M1/M2/M3 (Apple Silicon ARM64):**
Ejecuta esto en tu máquina para generar los binarios de Mac:
```powershell
$env:GOOS="darwin"; $env:GOARCH="arm64"
go build -ldflags="-s -w" -o api_mac ./cmd/api
go build -ldflags="-s -w" -o mock_backend_mac ./cmd/mock_backend
```
*Instrucciones para el QA en Mac:*
Debe abrir la terminal en la carpeta donde descargó los archivos, darles permiso de ejecución y correrlos:
```bash
chmod +x api_mac mock_backend_mac

# Pestaña 1: Levantar el Mock (puerto 8081)
./mock_backend_mac

# Pestaña 2: Levantar la API apuntando al Mock (puerto 8080 por defecto)
export FINNFLOW_URL="http://localhost:8081"
./api_mac
```

## Validación Final

Independientemente de la alternativa utilizada, el sistema estará operando localmente:
- **Swagger UI:** [http://localhost:8085/swagger-ui.html](http://localhost:8085/swagger-ui.html) *(o puerto 8080 según corresponda)*
- **Mock de FinnFlow:** [http://localhost:8081](http://localhost:8081)
