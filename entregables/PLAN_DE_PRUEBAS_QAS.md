# Plan de Pruebas (QAS)

## 1. Estrategia de Pruebas Automatizadas E2E (End-to-End)
Se creó un script automatizado `test_paridad.ps1` que lanza baterías de llamadas HTTP reales comparando el sistema Java con el sistema en Go, forzando la aserción de `Status Code Match` y `Body Match!` al 100%. Las variaciones de prueba garantizan la validación de la herencia del legado.

### Matrices de Pruebas Validadas E2E
| Caso de Prueba | Detalle | Resultado Físico Logrado |
| :--- | :--- | :--- |
| **TC-01** (ValidRequest) | Payload JSON completo. | 200 OK Idéntico. |
| **TC-02** (MissingField) | Se quita un nodo completo (participantes). Demuestra que Go salta las reglas de negocio al igual que Java. | Propagación del código y cuerpo devuelto por FinnFlow. |
| **TC-03** (InvalidValue_Monto) | Monto > Valor de la propiedad. Falla silenciosa omitida por legado. | Propagación del código y cuerpo devuelto por FinnFlow. |
| **TC-04** (Red Caída) | El Mock Backend se apaga, provocando fallo de conexión. | 500 y replicación de la cadena de excepción Java (I/O Error). |

## 2. Estrategia de Pruebas Unitarias y de Aislamiento
El microservicio en Go está testeado por capas (Clean Architecture), y todas las pruebas corren en el CI bajo `go test -race ./...`.

1.  **Capa de Transporte HTTP (`handler_test.go`):** Se inyectan peticiones falsas (mediante `httptest.NewRecorder`) hacia las rutas declaradas (ej. `/v1/bfcl/mortgage-loan/injections`). Evalúa que se retornen los códigos HTTP correctos (200 o status predeterminado por el error).
2.  **Capa de Casos de Uso (`workflow_test.go`):** Valida la orquestación. Usa un cliente de interfaz falsa (Mock) para devolver tokens inventados. Revisa que el empaquetamiento del response hacia el endpoint sea íntegro.
3.  **Capa de Integración a Terceros (`client_test.go`):** Levanta un servidor local hiperligero en memoria en milisegundos (`httptest.NewServer`). Envía peticiones HTTP, validando que el cliente extrae los JSONs del Token (`token_access`) adecuadamente.
