# Installs goremi on Windows: downloads the release binary to %LocalAppData%\Programs\goremi.
$ErrorActionPreference = 'Stop'

$releaseUrl = 'https://github.com/nhhthong/goremi/releases/latest/download'

# START: detect system
if ($env:PROCESSOR_ARCHITECTURE -ne 'AMD64') {
    [Console]::Error.WriteLine("goremi install: unsupported architecture: $env:PROCESSOR_ARCHITECTURE")
    exit 1
}
# END: detect system

# START: download and install
$dir = Join-Path (Join-Path $env:LocalAppData 'Programs') 'goremi'
New-Item -ItemType Directory -Force -Path $dir | Out-Null
Invoke-WebRequest -Uri "$releaseUrl/goremi_windows_amd64.exe" -OutFile (Join-Path $dir 'goremi.exe')
# END: download and install

# START: check dependencies
foreach ($tool in 'mpv', 'yt-dlp') {
    if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
        Write-Host "$tool not found"
    }
}
# END: check dependencies
