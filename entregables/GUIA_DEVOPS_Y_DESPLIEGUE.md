# Guía de DevOps y Despliegue (Producción)

## 1. Empaquetado en Contenedores (Docker)
El proyecto utiliza un patrón de construcción Multi-Etapa (Multi-Stage) altamente agresivo y optimizado.

### Construcción
```bash
docker build -t mortgage-loan-injection-go:latest .
```
1.  **Stage 1 (Builder):** Usa `golang:1.23-alpine`. Descarga dependencias y compila el código máquina estáticamente anulando CGO (`CGO_ENABLED=0`). Utiliza los flags `-ldflags="-w -s"` para eliminar símbolos de depuración, recortando el tamaño a la mitad.
2.  **Stage 2 (Scratch):** Monta la imagen `scratch` de Docker (0 MB). Inserta los certificados SSL raíces extraídos del builder y el binario ultra-ligero.

## 2. Integración Continua (CI Pipeline)
Se configuró **GitHub Actions** (`.github/workflows/ci.yml`). Cada vez que se integre código (Push o PR) a `main` o `develop`:
*   Levanta un runner `ubuntu-latest`.
*   Configura el entorno de Go 1.23 con estrategias de caché agresivas.
*   Verifica formateos e importaciones huérfanas (`go mod tidy`, `go vet`).
*   Corre la batería de pruebas unitarias inyectando el detector de condiciones de carrera (`-race`).
*   Exporta y empaqueta el artefacto de Cobertura (`coverage.txt`).

## 3. Despliegue en Kubernetes (CD)
Se generaron manifiestos nativos (`k8s/deployment.yaml` y `k8s/service.yaml`).
*   **Resource Limits:** Maximizando la rentabilidad, los límites se establecieron en *CPU 250m* y *Memory 128Mi* (Mínimos históricos inalcanzables en JVM).
*   **Inyección de Secretos:** Integración nativa con `SecretKeyRef` para montar el `FINNFLOW_USER` y `FINNFLOW_PASS` de forma segura.
*   **Probes de Salud:** `livenessProbe` y `readinessProbe` apuntando al endpoint ultra-rápido de memoria inyectada (`/api-docs`).
