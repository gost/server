@echo off

if exist ..\bin rd /s /q ..\bin
mkdir ..\bin\linuxARMv6
mkdir ..\bin\linuxARMv7

SET GOOS=linux
SET GOARCH=arm
SET GOARM=6
go build -o ..\bin\linuxARMv6\gost ../src/main.go
echo "Built application for Linux/ARMv6"
SET GOARM=7
go build -o ..\bin\linuxARMv7\gost ../src/main.go
echo "Built application for Linux/ARMv7"
