# Installs goremi on Windows: downloads the release binary to %LocalAppData%\Programs\goremi.
param([switch]$Yes)
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
$missing = @()
foreach ($tool in 'mpv', 'yt-dlp') {
    if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
        Write-Host "$tool not found"
        $missing += $tool
    }
}
# END: check dependencies

# START: offer dependencies
if ($missing.Count -gt 0) {
    $haveWinget = [bool](Get-Command winget -ErrorAction SilentlyContinue)
    $haveScoop = [bool](Get-Command scoop -ErrorAction SilentlyContinue)
    $ytUrl = 'https://github.com/yt-dlp/yt-dlp/releases/latest/download'
    $ytTarget = Join-Path $dir 'yt-dlp.exe'
    Write-Host 'These commands install what is missing:'
    if ($missing -contains 'mpv') {
        # With neither winget nor scoop both commands are shown, and none is run.
        if ($haveWinget -or -not $haveScoop) { Write-Host '  winget install --id shinchiro.mpv' }
        if ($haveScoop -or -not $haveWinget) { Write-Host '  scoop bucket add extras; scoop install mpv' }
    }
    if ($missing -contains 'yt-dlp') {
        Write-Host "  Invoke-WebRequest $ytUrl/yt-dlp.exe -OutFile $ytTarget"
    }
    $answer = 'y'
    if (-not $Yes) {
        try { $answer = Read-Host ('Install ' + ($missing -join ' and ') + '? [y/N]') } catch { $answer = '' }
        if ($answer -notmatch '^(y|yes)$') { $answer = '' }
    }
    if ($answer -ne '') {
        if ($missing -contains 'mpv' -and ($haveWinget -or $haveScoop)) {
            $global:LASTEXITCODE = 0
            if ($haveWinget) {
                & winget install --id shinchiro.mpv
            } else {
                & scoop bucket add extras
                if ($LASTEXITCODE -eq 0) { & scoop install mpv }
            }
            if ($LASTEXITCODE -ne 0) { [Console]::Error.WriteLine('goremi install: could not install mpv') }
        }
        if ($missing -contains 'yt-dlp') {
            $partial = "$ytTarget.download"
            $sums = Join-Path $dir 'SHA2-256SUMS.download'
            try {
                Invoke-WebRequest -Uri "$ytUrl/yt-dlp.exe" -OutFile $partial
                Invoke-WebRequest -Uri "$ytUrl/SHA2-256SUMS" -OutFile $sums
                $want = ((Get-Content -LiteralPath $sums | Where-Object { $_ -match '\s\*?yt-dlp\.exe$' } | Select-Object -First 1) -split '\s+')[0]
                $got = (Get-FileHash -Algorithm SHA256 -LiteralPath $partial).Hash
                if (-not $want -or $want -ne $got) { throw 'the checksum of yt-dlp.exe does not match SHA2-256SUMS' }
                Move-Item -Force -LiteralPath $partial -Destination $ytTarget
            } catch {
                [Console]::Error.WriteLine("goremi install: could not install yt-dlp: $($_.Exception.Message)")
            } finally {
                Remove-Item -Force -LiteralPath $partial, $sums -ErrorAction SilentlyContinue
            }
        }
    }
}
# END: offer dependencies

# START: check PATH
if (-not (($env:Path -split [IO.Path]::PathSeparator) -contains $dir)) {
    Write-Host "$dir is not in your PATH; add it to run goremi from any folder."
}
# END: check PATH
