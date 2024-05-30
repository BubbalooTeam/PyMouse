@echo off

setlocal
for %%i in (%~dp0) do set "parent_dir=%%~dpi"
endlocal & set "parent_dir=%parent_dir%"

:: Check if the user already has administrator privileges
net session >nul 2>&1
if %errorlevel% neq 0 (
    :: If not, run the batch script as an administrator
    echo Running as administrator...
    runas /user:Administrator 
    start "%parent_dir%\PyMouse\resources\windows\utils\scripts\win_installer_pymouse.bat"
    echo PyMouse installation is complete!
    exit
)

:: Only start the installer if the user is already an administrator
start "%parent_dir%\PyMouse\resources\windows\utils\scripts\win_installer_pymouse.bat"
echo PyMouse installation is complete!
pause

