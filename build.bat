@echo off
setlocal

REM 设置 CGO 关闭，保证静态编译
set CGO_ENABLED=0

set VERSION=0.1.5

REM ================================
REM Windows 64-bit
REM ================================
REM set GOOS=windows
REM set GOARCH=amd64
REM go build -o bin\weixinagent-win-amd64-%VERSION%.exe .

REM ================================
REM Linux 64-bit
REM ================================
REM set GOOS=linux
REM set GOARCH=amd64
REM go build -o bin\weixinagent-linux-amd64-%VERSION%  .

REM ================================
REM Linux 64-bit
REM ================================
REM set GOOS=linux
REM set GOARCH=arm64
REM go build -o bin\weixinagent-linux-arm64-%VERSION% .

REM ================================
REM macOS Intel (amd64)
REM ================================
set GOOS=darwin
set GOARCH=amd64
go build -o bin\weixinagent-macos-amd64-%VERSION% .

REM ================================
REM macOS ARM (M1/M2, arm64)
REM ================================
REM set GOOS=darwin
REM set GOARCH=arm64
REM go build -o bin\weixinagent-macos-arm64-%VERSION% .

echo.
echo weixinagent BUILD successfully
echo.

endlocal
