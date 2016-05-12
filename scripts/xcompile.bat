@echo off

if exist ..\bin rd /s /q ..\bin
mkdir ..\bin\win32
mkdir ..\bin\win64
mkdir ..\bin\linux32
mkdir ..\bin\linux64
mkdir ..\bin\linuxARMv6
mkdir ..\bin\linuxARMv7
mkdir ..\bin\darwin32
mkdir ..\bin\darwin64

SET GOOS=windows
SET GOARCH=386
go build -o ..\win32/gost.exe ../src/main.go
echo "Built application for Windows/386"
SET GOARCH=amd64
go build -o ..\win64/gost.exe ../src/main.go
echo "Built application for Windows/amd64"

SET GOOS=linux
SET GOARCH=386
go build -o ..\linux32/gost ../src/main.go
echo "Built application for Linux/386"
SET GOARCH=amd64
go build -o ..\linux64/gost ../src/main.go
echo "Built application for Linux/amd64"
SET GOARCH=arm
SET GOARM=6
go build -o ..\linuxARMv6/gost ../src/main.go
echo "Built application for Linux/ARMv6"
SET GOARM=7
go build -o ..\linuxARMv7/gost ../src/main.go
echo "Built application for Linux/ARMv7"

SET GOOS=darwin
SET GOARCH=386
go build -o ..\darwin32/gost ../src/main.go
echo "Built application for Darwin/386"
SET GOARCH=amd64
go build -o ..\darwin64/gost ../src/main.go
echo "Built application for Darwin/amd64"