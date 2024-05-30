@echo off
echo Updating system to install BOT, please wait...

:: Check if Chocolatey is installed
where choco > nul 2>&1
if %errorlevel% neq 0 (
    echo Chocolatey not found, installing Chocolatey...
    @powershell -NoProfile -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))"
    echo Chocolatey installed successfully!
)

:: Check if Python is installed
where python > nul 2>&1
if %errorlevel% neq 0 (
    echo Python not found, installing Python...
    choco install python -y
    echo Python and tools installed successfully!
)

echo Updates of Termux are terminated! Initializing Install PyMouse Tool...