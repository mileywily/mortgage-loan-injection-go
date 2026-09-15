# test_paridad.ps1

Write-Host "Iniciando validación de paridad E2E (Java vs Go)..."

# URLs
$JAVA_URL = "http://localhost:8080/v1/bfcl/mortgage-loan/injections"
$GO_URL = "http://localhost:8085/v1/bfcl/mortgage-loan/injections"

# Test cases
$payloads = @{
    "ValidRequest" = '{"datos_credito": {"NumeroSolicitud": 12345, "AntiguedadVivienda": 5, "EjecutivoComercial": "Juan", "Producto": 1, "Objetivo": 1, "Destino": 1, "MontoAprobado": 1000.0, "ValorPropiedad": 2000.0, "FechaAprobacion": "2026-01-01", "ValorContado": 500.0, "Plazo1": 15, "MesesGracia": 0, "Tasa1": 4.5, "Spread1": 1.2}, "participantes": [{"Rut": "12345678-9", "Nombre": "Juan", "Apellido": "Perez", "Rol": "Titular"}]}'
    "MissingField" = '{"datos_credito": {"AntiguedadVivienda": 5}, "participantes": []}'
    "InvalidValue_Monto" = '{"datos_credito": {"NumeroSolicitud": 12345, "AntiguedadVivienda": 5, "EjecutivoComercial": "Juan", "Producto": 1, "Objetivo": 1, "Destino": 1, "MontoAprobado": 3000.0, "ValorPropiedad": 2000.0, "FechaAprobacion": "2026-01-01", "ValorContado": 500.0, "Plazo1": 15, "MesesGracia": 0, "Tasa1": 4.5, "Spread1": 1.2}, "participantes": [{"Rut": "12345678-9"}]}'
    "InvalidFormat_Plazo" = '{"datos_credito": {"NumeroSolicitud": 12345, "AntiguedadVivienda": 5, "EjecutivoComercial": "Juan", "Producto": 1, "Objetivo": 1, "Destino": 1, "MontoAprobado": 1000.0, "ValorPropiedad": 2000.0, "FechaAprobacion": "2026-01-01", "ValorContado": 500.0, "Plazo1": 17, "MesesGracia": 0, "Tasa1": 4.5, "Spread1": 1.2}, "participantes": [{"Rut": "12345678-9"}]}'
}

foreach ($testName in $payloads.Keys) {
    $payload = $payloads[$testName]
    Write-Host "`nEjecutando Test: $testName"

    $javaResult = ""
    $javaStatus = ""
    $goResult = ""
    $goStatus = ""

    # Call Java using curl to avoid Invoke-WebRequest stream bug
    $javaRaw = curl.exe -s -w "\n%{http_code}" -X POST $JAVA_URL -H "Content-Type: application/json" -d $payload
    $javaLines = $javaRaw -split "`n"
    $javaStatus = [int]($javaLines[-1])
    $javaResult = ($javaLines[0..($javaLines.Length-2)] -join "`n").Trim()

    # Call Go using curl
    $goRaw = curl.exe -s -w "\n%{http_code}" -X POST $GO_URL -H "Content-Type: application/json" -d $payload
    $goLines = $goRaw -split "`n"
    $goStatus = [int]($goLines[-1])
    $goResult = ($goLines[0..($goLines.Length-2)] -join "`n").Trim()

    if ($javaStatus -eq $goStatus) {
        Write-Host "Status Code Match: $javaStatus" -ForegroundColor Green
    } else {
        Write-Host "Status Code MISMATCH: Java=$javaStatus vs Go=$goStatus" -ForegroundColor Red
    }

    if ($javaResult.Trim() -eq $goResult.Trim()) {
        Write-Host "Body Match!" -ForegroundColor Green
    } else {
        Write-Host "Body MISMATCH!" -ForegroundColor Red
        Write-Host "Java Body: $javaResult"
        Write-Host "Go Body: $goResult"
    }
}
