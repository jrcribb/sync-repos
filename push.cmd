@echo off
REM Script para hacer push con la versión correcta
REM Versión 23.12.2
GOTO :MAIN
if "%1"=="" (
    set 
)
:BUILDNUMBER
if not exist "buildnumber.cfg" (
    echo 0 > buildnumber.cfg
)
set /p BuildNumber=<buildnumber.cfg
CALL :ANIO
set buildnumber=v.%year%.%month%/%BuildNumber%
GOTO :EOF

:ANIO
set year=%date:~-2,2%
set month=%date:~-7,2%
GOTO :EOF

:MAIN
CALL :BUILDNUMBER
git add -A
git commit -m %buildnumber%
git push