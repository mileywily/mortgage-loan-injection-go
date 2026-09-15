# start_dev.ps1
# Script para iniciar el microservicio en Go y el Mock (Dummy) Backend localmente sin usar Docker.

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Iniciando Entorno de Pruebas (Sin Docker)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# 1. Matar procesos anteriores si quedaron pegados
$portsToKill = @(8081, 8085)
foreach ($port in $portsToKill) {
    $pidToKill = (Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess | Select-Object -First 1)
    if ($pidToKill) {
        Write-Host "Cerrando proceso antiguo en puerto $port (PID: $pidToKill)..." -ForegroundColor Yellow
        Stop-Process -Id $pidToKill -Force -ErrorAction SilentlyContinue
    }
}

# 2. Levantar el Mock Backend (Dummy FinnFlow)
Write-Host "`n[1/2] Levantando Mock Backend (Dummy FinnFlow) en el puerto 8081..." -ForegroundColor Green
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "./cmd/mock_backend"

# Darle 2 segundos para que arranque
Start-Sleep -Seconds 2

# 3. Levantar la API en Go apuntando al Mock
Write-Host "`n[2/2] Levantando API Principal en Go en el puerto 8085..." -ForegroundColor Green
$env:PORT = "8085"
$env:FINNFLOW_URL = "http://localhost:8081"
$env:FINNFLOW_USER = "testuser"
$env:FINNFLOW_PASS = "testpass"

Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "./cmd/api"

Write-Host "`n==========================================" -ForegroundColor Cyan
Write-Host "¡Todo en línea y funcionando!" -ForegroundColor Green
Write-Host "Swagger UI : http://localhost:8085/swagger-ui.html"
Write-Host "Endpoint   : http://localhost:8085/v1/bfcl/mortgage-loan/injections"
Write-Host "Mock       : http://localhost:8081"
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "NOTA: Puedes probar la inyección usando Postman hacia el puerto 8085."
Write-Host "Para apagar los servicios, presiona Ctrl+C o cierra esta ventana."
