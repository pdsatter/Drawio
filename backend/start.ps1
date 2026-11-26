param($Env:PORT="8080", $Cmd="/cmd/server/main.go")

# Ensure Go path is in PATH (adjust for Windows)
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

try {
    $proc = Start-Process -FilePath go -ArgumentList "run", "$Cmd" -Wait -PassThru
} catch {
    Write-Host "Failed to start Go server:" $_.Exception.Message
    exit 1
}
