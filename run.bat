@echo off
SETLOCAL

SET DBNAME=eventboarddb
SET DBUSER=postgres
SET CACHE=false
SET PRODUCTION=false
SET DBPASS=tom
SET DBPORT=5432

cd cmd\web

go build -o ../../nastenka_udalosti.exe

IF NOT %ERRORLEVEL% == 0 (
    echo Build failed, exiting...
    goto :eof
)

cd ..\..

nastenka_udalosti.exe -dbname=%DBNAME% -dbuser=%DBUSER% -cache=%CACHE% -production=%PRODUCTION% -dbpass=%DBPASS% -dbport=%DBPORT%

ENDLOCAL
