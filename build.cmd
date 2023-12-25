@echo off
REM Script para generar un número de compilación para usar en la versión generada del sistema.
REM Versión 23.12.3
GOTO :MAIN

:BUILDNUMBER
if not exist "buildnumber.cfg" (
    echo 0 > buildnumber.cfg
)
set /p BuildNumber=<buildnumber.cfg
set /a build=%BuildNumber%+1
echo %build% > buildnumber.cfg
CALL :ANIO
set buildnumber=v.%year%.%month%/%build%
GOTO :EOF

:ANIO
set year=%date:~-2,2%
set month=%date:~-7,2%
GOTO :EOF

:MAIN
CALL :BUILDNUMBER
go build -ldflags "-X main.version=%buildnumber%"