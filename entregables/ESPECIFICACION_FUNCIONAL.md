# Especificación Funcional y No Funcional

## 1. Especificaciones Funcionales

### Caso de Uso Principal (CU-01)
*   **Nombre:** Inyección Transparente de Solicitud de Crédito a FinnFlow.
*   **Descripción:** El sistema actúa como middleware. Recibe solicitudes de crédito, obtiene un token de seguridad, y reenvía (inyecta) los datos a Salesforce/FinnFlow.
*   **Reglas de Negocio (El Descubrimiento de la Deuda Técnica):** Aunque en el código legado existían funciones teóricas para validar topes de montos, plazos múltiplos de 5, rangos de tasas y existencias de participantes, estas **nunca eran invocadas por el controlador principal**. En consecuencia, la directiva funcional para este reemplazo es **omitir cualquier validación de negocio propia**, aceptando cualquier JSON válido y delegando la decisión 100% al sistema final (FinnFlow).

### Gestión de Errores (Error Handling)
El microservicio propaga el estado del proveedor de forma directa:
1.  **Errores Lógicos (HTTP 400):** Si FinnFlow rechaza la solicitud, el MS captura la respuesta, la parsea y la retransmite con el mismo cuerpo al consumidor original.
2.  **Errores de Autenticación (HTTP 401):** Retorna `500 Internal Server Error` detallando la falla de autenticación, mapeando la respuesta exacta del legado.
3.  **Errores de Red (TimeOut):** Retorna `500 Internal Server Error` y replica el mensaje de excepción crudo de I/O, en estricto cumplimiento con la paridad exigida.

## 2. Especificaciones No Funcionales

*   **Rendimiento y Latencia:** Al estar construido en Go 1.23 nativo (sin recolección de basura pausada estilo JVM), la sobrecarga computacional del proxy se reduce al mínimo, sumando < 2ms de latencia al Round-Trip-Time (RTT) contra FinnFlow.
*   **Huella de Memoria (Footprint):** El consumo en reposo (*Idle*) es inferior a **5 MB RAM**, comparado con los > 250 MB habituales del contenedor de Spring Boot previo.
*   **Tamaño del Desplegable (Storage):** El ejecutable compilado estáticamente pesa < 15 MB, a diferencia de los JARs pesados empaquetados junto al JRE.
*   **Resiliencia SSL (Ambientes Controlados):** Acepta el bypass condicional de certificados autofirmados mediante la variable `INSECURE_SKIP_VERIFY=true`.
