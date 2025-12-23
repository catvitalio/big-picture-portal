@echo off
setlocal

if "%1"=="clean" goto clean
if "%1"=="install-tools" goto install-tools
if "%1"=="rebuild" goto rebuild
goto build

:build
echo Building Big Picture Portal...

REM Generate resource file from manifest if not exists or manifest is newer
if not exist rsrc.syso (
    call :generate-rsrc
) else (
    for %%i in (app.manifest) do set manifest_time=%%~ti
    for %%i in (rsrc.syso) do set rsrc_time=%%~ti
    if "!manifest_time!" GTR "!rsrc_time!" (
        call :generate-rsrc
    )
)

REM Build the application
go build -ldflags="-H windowsgui" -o BigPicturePortal.exe
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Build complete: BigPicturePortal.exe
goto end

:generate-rsrc
echo Generating resource file...

REM Check if icon exists and add it to rsrc command
if exist assets\icon.ico (
    rsrc -manifest app.manifest -ico assets\icon.ico -o rsrc.syso 2>nul
    if errorlevel 1 (
        echo rsrc tool not found, installing...
        go install github.com/akavel/rsrc@latest
        rsrc -manifest app.manifest -ico assets\icon.ico -o rsrc.syso
    )
) else (
    rsrc -manifest app.manifest -o rsrc.syso 2>nul
    if errorlevel 1 (
        echo rsrc tool not found, installing...
        go install github.com/akavel/rsrc@latest
        rsrc -manifest app.manifest -o rsrc.syso
    )
)
goto :eof

:clean
echo Cleaning build artifacts...
if exist BigPicturePortal.exe del BigPicturePortal.exe
if exist rsrc.syso del rsrc.syso
echo Clean complete.
goto end

:install-tools
echo Installing build tools...
go install github.com/akavel/rsrc@latest
echo Tools installed.
goto end

:rebuild
call :clean
call :build
goto end

:end
endlocal