@echo off
SETLOCAL

@REM Nutno zadat jméno databáse
SET DBNAME=eventboarddb
@REM Nutno zadat jméno - (default - postgres)
SET DBUSER=postgres
SET CACHE=true
SET PRODUCTION=true
@REM Nutno zadat heslo
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
