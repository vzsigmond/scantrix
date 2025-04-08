# install.ps1 - Scantrix Windows Installer

$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$binary = "scantrix-windows-$arch.exe"
$url = "https://github.com/vzsigmond/scantrix/releases/latest/download/$binary"
$tempOut = "$env:TEMP\scantrix.exe"
$installDir = "$env:ProgramData\scantrix"
$installPath = "$installDir\scantrix.exe"

Write-Host "⬇️ Downloading Scantrix binary from $url..."
Invoke-WebRequest -Uri $url -OutFile $tempOut

Write-Host "📁 Installing to $installPath"
New-Item -Path $installDir -ItemType Directory -Force | Out-Null
Move-Item -Path $tempOut -Destination $installPath -Force

# Add to system PATH if not already present
$envPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($envPath -notlike "*scantrix*") {
    Write-Host "🔧 Adding $installDir to system PATH"
    [Environment]::SetEnvironmentVariable("Path", "$envPath;$installDir", "Machine")
}

Write-Host "`n✅ Scantrix installed successfully!"
Write-Host "👉 Restart your terminal or run 'scantrix --help'"
