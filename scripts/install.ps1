# Keyarch installer for Windows
# Usage: irm https://raw.githubusercontent.com/vinayydv3695/keyarch/main/scripts/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "vinayydv3695/keyarch"
$BinaryName = "keyarch.exe"

# Determine install directory
$InstallDir = "$env:LOCALAPPDATA\Programs\keyarch"

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] " -ForegroundColor Red -NoNewline
    Write-Host $Message
    exit 1
}

function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    try {
        $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $Release.tag_name
    } catch {
        Write-Error "Failed to get latest version: $_"
    }
}

function Get-Architecture {
    $Arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")
    switch ($Arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { Write-Error "Unsupported architecture: $Arch" }
    }
}

function Install-Keyarch {
    param(
        [string]$Version,
        [string]$Arch
    )
    
    $DownloadUrl = "https://github.com/$Repo/releases/download/$Version/keyarch-windows-$Arch.exe"
    
    Write-Info "Downloading from: $DownloadUrl"
    
    # Create install directory
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    $DestPath = Join-Path $InstallDir $BinaryName
    
    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $DestPath -UseBasicParsing
    } catch {
        Write-Error "Failed to download: $_"
    }
    
    Write-Info "Installed to: $DestPath"
    
    return $DestPath
}

function Add-ToPath {
    param([string]$Directory)
    
    $CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    if ($CurrentPath -notlike "*$Directory*") {
        Write-Info "Adding $Directory to PATH..."
        $NewPath = "$CurrentPath;$Directory"
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        $env:Path = "$env:Path;$Directory"
        Write-Info "Added to PATH. You may need to restart your terminal."
    } else {
        Write-Info "$Directory is already in PATH"
    }
}

function Verify-Installation {
    try {
        $Version = & keyarch --help 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Info "Installation successful!"
            Write-Host ""
            Write-Host "Run 'keyarch' to start typing!" -ForegroundColor Cyan
            return $true
        }
    } catch {
        # Ignore errors
    }
    return $false
}

function Main {
    Write-Host ""
    Write-Host "  Keyarch Installer for Windows" -ForegroundColor Cyan
    Write-Host "  Feature-rich TUI typing test"
    Write-Host ""
    
    $Version = Get-LatestVersion
    Write-Info "Latest version: $Version"
    
    $Arch = Get-Architecture
    Write-Info "Detected architecture: $Arch"
    
    $InstalledPath = Install-Keyarch -Version $Version -Arch $Arch
    
    Add-ToPath -Directory $InstallDir
    
    Write-Host ""
    
    if (-not (Verify-Installation)) {
        Write-Warn "Please restart your terminal and run 'keyarch' to start!"
    }
    
    Write-Host ""
}

Main
