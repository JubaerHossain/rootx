$ErrorActionPreference = "Stop"

$BinDir = "C:\Program Files\rootx"
$BinName = "rootx.exe"

if (-not (Test-Path -Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir
}

Copy-Item -Path $BinName -Destination $BinDir

$Path = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::Machine)
if ($Path -notlike "*$BinDir*") {
    [System.Environment]::SetEnvironmentVariable("Path", "$Path;$BinDir", [System.EnvironmentVariableTarget]::Machine)
}

Write-Output "rootx installed successfully to $BinDir"
