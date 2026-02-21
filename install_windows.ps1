# Check if running as administrator
if (-Not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) {
    Write-Host "Please run this script as Administrator" -ForegroundColor Red
    exit 1
}

Write-Host "Building the Binary"
go build -o agai.exe .

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}

Write-Host "Installing ..."
$system32 = [System.Environment]::GetFolderPath([System.EnvironmentSpecialFolder]::System)
Move-Item -Path "agai.exe" -Destination "$system32\agai.exe" -Force

Write-Host "Installation Done" -ForegroundColor Green
agai -h
