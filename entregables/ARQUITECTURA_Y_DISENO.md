# Arquitectura y Diseño

## 1. Patrón Arquitectónico: Clean Architecture & Go-Kit
El microservicio ha sido reescrito utilizando los principios de **Clean Architecture**, aislando estrictamente la lógica del transporte de red de las reglas de negocio. Para estructurar este patrón de manera estandarizada, se implementó la abstracción inspirada en **Go-Kit**:

*   **Transport Layer (`internal/handler`):** Maneja puramente el protocolo HTTP, los ruteos y la decodificación de JSON (request/response).
*   **Endpoint Layer (`internal/endpoint`):** Interfaz que adapta la solicitud HTTP entrante a un formato agnóstico de red para ser procesado por el core.
*   **Service/Workflow Layer (`internal/workflows`):** Contiene la orquestación del negocio (obtención del token y llamada de inyección).
*   **Activities/Client Layer (`internal/activities`):** Capa de salida (Egress) que interactúa con la infraestructura de terceros (Salesforce / FinnFlow).

## 2. Paradigma de Gestión de Estado: Stateless & Zero-Cache
Tras analizar el microservicio legado (Spring Boot), se dictaminó que la arquitectura real operaba como un *Reverse Proxy Pass-Through*.
*   **Sin Estado (Stateless):** El sistema en Go no guarda persistencia en memoria ni en disco de los payloads transaccionados.
*   **Zero-Cache:** El Token de Autenticación hacia FinnFlow no se guarda en memoria. Cada petición entrante gatilla de forma síncrona una petición de autenticación (`POST /api/token`) seguida de la inyección (`POST /api/inyeccion-salesforce`). Esto facilita un escalamiento horizontal infinito en Kubernetes sin fricciones por desincronización de cachés.

## 3. Observabilidad Nativa (JSON Logging Parity)
Para no romper la ingesta de logs en los colectores empresariales (Datadog, ELK, Splunk) que esperaban el formato de la librería `logback-json-classic` de Java, se configuró el motor `log/slog` de Go (nativo a partir de 1.21).
Mediante un interceptor `ReplaceAttr`, las llaves nativas de Go (`"time"`, `"msg"`) son renombradas al vuelo a `"timestamp"` y `"message"`, logrando **paridad extrema byte-for-byte** en los logs estándar.

## 4. Portabilidad y Auto-Contención
A través de la directiva `//go:embed`, se embebieron en memoria los artefactos de documentación (`swagger.json` y `swagger-ui.html`). El microservicio se compila en un único binario ejecutable (`api.exe` o `api`) que no depende del sistema de archivos host para proveer su interfaz.
