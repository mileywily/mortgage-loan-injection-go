# Guía QA y Despliegue Local (Sin Docker)

Para equipos de pruebas o evaluadores funcionales que no dispongan de contenedores (Docker) instalados ni permisos de virtualización, el ecosistema puede correr nativamente.

## Componentes Disponibles
El proyecto entrega 2 componentes:
1.  **API Core (`api`)**: El microservicio de inyección en sí mismo.
2.  **Mock Backend (`mock_backend`)**: Un simulador hiper-ligero (Dummy) que se hace pasar por el sistema corporativo FinnFlow para que la API responda exitosamente.

## Opción 1: Compilación Rápida Local (Usuarios con Go Instalado)
Si dispones del lenguaje Go instalado:
*   En Windows: Simplemente haz clic derecho y ejecuta con PowerShell el script provisto: `.\scripts\start_dev.ps1`.
*   Esto levantará el Mock en el puerto `8081` y la API apuntando a él en el puerto `8085`. 

## Opción 2: Ejecutables Pre-Compilados "Double-Click" (Zero Dependencies)
Dada la maravillosa habilidad de **Compilación Cruzada** de Go, el desarrollador te compartirá dos archivos ejecutables nativos hechos a medida para tu computadora. ¡Ni siquiera necesitas instalar Go!

### Si usas Windows
Recibirás: `api.exe` y `mock_backend.exe`.
1.  Doble clic a `mock_backend.exe`.
2.  Doble clic a `api.exe`.
¡Listo!

### Si usas Mac (M1/M2/M3 - Apple Silicon ARM64)
Recibirás: `api_mac` y `mock_backend_mac`.
1.  Abre la terminal. Otorga permisos de ejecución: `chmod +x api_mac mock_backend_mac`
2.  En una pestaña de la terminal, arranca el simulador: `./mock_backend_mac`
3.  En otra pestaña, indícale a la API dónde está el simulador y arráncala: 
    ```bash
    export FINNFLOW_URL="http://localhost:8081"
    ./api_mac
    ```

## Acceso Final
Ambas opciones pondrán en línea tu Swagger para realizar pruebas directas:
*   **Interfaz Gráfica:** `http://localhost:8085/swagger-ui.html`
*   **Punto de Inyección:** `http://localhost:8085/v1/bfcl/mortgage-loan/injections`
