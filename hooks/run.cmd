@echo off
setlocal

set "ARCHITECTURE=%PROCESSOR_ARCHITECTURE%"
if defined PROCESSOR_ARCHITEW6432 set "ARCHITECTURE=%PROCESSOR_ARCHITEW6432%"

if /i "%ARCHITECTURE%"=="AMD64" set "TARGET_ARCH=amd64"
if /i "%ARCHITECTURE%"=="x86_64" set "TARGET_ARCH=amd64"
if /i "%ARCHITECTURE%"=="ARM64" set "TARGET_ARCH=arm64"
if /i "%ARCHITECTURE%"=="aarch64" set "TARGET_ARCH=arm64"

if not defined TARGET_ARCH (
  >&2 echo codex-next-prompt: unsupported architecture: %ARCHITECTURE%
  exit /b 0
)

"%~dp0..\bin\windows-%TARGET_ARCH%\codex-next-prompt.exe" %*
exit /b %ERRORLEVEL%
