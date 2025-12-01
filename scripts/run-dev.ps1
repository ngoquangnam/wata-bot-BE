# PowerShell script to run in development environment
# Usage: .\scripts\run-dev.ps1

Write-Host "Starting Wata Bot Backend in DEV mode..." -ForegroundColor Green

# Load dev environment variables
if (Test-Path .env.dev) {
    Get-Content .env.dev | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
    Write-Host "Loaded .env.dev file" -ForegroundColor Yellow
} else {
    Write-Host "Warning: .env.dev file not found, using config file only" -ForegroundColor Yellow
}

# Run the application with dev config
go run wata-bot.go -f etc/wata-bot-api.dev.yaml

