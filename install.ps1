# Installs goremi on Windows: downloads the release binary to %LocalAppData%\Programs\goremi.
$ErrorActionPreference = 'Stop'
# Windows PowerShell 5.1 downloads very slowly while the progress bar shows.
$ProgressPreference = 'SilentlyContinue'

$releaseUrl = 'https://github.com/nhhthong/goremi/releases/latest/download'

# START: detect system
# A 32-bit PowerShell on 64-bit Windows reports x86 in PROCESSOR_ARCHITECTURE; PROCESSOR_ARCHITEW6432 then holds the real one.
$arch = $env:PROCESSOR_ARCHITEW6432
if ([string]::IsNullOrEmpty($arch)) { $arch = $env:PROCESSOR_ARCHITECTURE }
if ($arch -ne 'AMD64') {
    [Console]::Error.WriteLine("goremi install: unsupported architecture: $arch (PROCESSOR_ARCHITECTURE=$env:PROCESSOR_ARCHITECTURE, PROCESSOR_ARCHITEW6432=$env:PROCESSOR_ARCHITEW6432)")
    exit 1
}
# END: detect system

# START: download and install
$dir = Join-Path (Join-Path $env:LocalAppData 'Programs') 'goremi'
New-Item -ItemType Directory -Force -Path $dir | Out-Null
$target = Join-Path $dir 'goremi.exe'
$partial = "$target.download"
try {
    Invoke-WebRequest -Uri "$releaseUrl/goremi_windows_amd64.exe" -OutFile $partial
    Move-Item -Force -LiteralPath $partial -Destination $target
} catch {
    Remove-Item -Force -LiteralPath $partial -ErrorAction SilentlyContinue
    [Console]::Error.WriteLine("goremi install: could not install goremi.exe: $($_.Exception.Message) Close goremi if it is running, then run this script again.")
    exit 1
}
# END: download and install

# START: check dependencies
$hints = @{
    'mpv'    = 'winget install --id shinchiro.mpv, or: scoop bucket add extras; scoop install mpv'
    'yt-dlp' = 'winget install --id yt-dlp.yt-dlp, or: scoop install yt-dlp'
}
foreach ($tool in 'mpv', 'yt-dlp') {
    if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
        Write-Host "$tool not found (install it with: $($hints[$tool]))"
    }
}
# END: check dependencies

# START: check PATH
if (-not (($env:Path -split [IO.Path]::PathSeparator) -contains $dir)) {
    Write-Host "$dir is not in your PATH; add it to run goremi from any folder."
}
# END: check PATH
