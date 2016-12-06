@echo off

if exist ..\bin rd /s /q ..\bin
mkdir ..\bin\win64
mkdir ..\bin\linux64
mkdir ..\bin\darwin64

SET GOOS=windows
SET GOARCH=amd64
go build -o ..\bin\win64\gost.exe ../src/main.go
echo "Built application for Windows/amd64"
copy ..\src\config.yaml ..\bin\win64\config.yaml
xcopy ..\src\client\*.* ..\bin\win64\client\ /S

SET GOOS=linux
SET GOARCH=amd64
go build -o ..\bin\linux64\gost ../src/main.go
echo "Built application for Linux/amd64"
copy ..\src\config.yaml ..\bin\linux64\config.yaml
xcopy ..\src\client\*.* ..\bin\linux64\client\ /S

SET GOOS=darwin
SET GOARCH=amd64
go build -o ..\bin\darwin64\gost ../src/main.go
echo "Built application for Darwin/amd64"
copy ..\src\config.yaml ..\bin\darwin64\config.yaml
xcopy ..\src\client\*.* ..\bin\darwin64\client\ /S
