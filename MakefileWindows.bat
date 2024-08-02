@echo off
setlocal

set BIN_DIR=bin
set OUTPUT_FILE=%BIN_DIR%\E-commerce_REST-API.exe
set CMD_FILE=cmd\main.go
set MIGRATION_DIR=cmd\migrate\migrations
set MIGRATE_CMD=C:\Users\ricar\Documents\GitHub\E-commerce_REST-API\cmd\migrate\migrate.exe

if "%1" == "build" (
    if not exist %BIN_DIR% (
        mkdir %BIN_DIR%
    )
    echo Building Go application...
    go build -o %OUTPUT_FILE% %CMD_FILE%
    if exist %OUTPUT_FILE% (
        echo Build successful.
    ) else (
        echo Build failed.
        exit /b 1
    )
)

if "%1" == "test" (
    echo Running Go tests...
    go test -v ./...
)

if "%1" == "run" (
    echo Running Go application...
    if exist %OUTPUT_FILE% (
        go build -o %OUTPUT_FILE% %CMD_FILE%
        %OUTPUT_FILE%
    ) else (
        echo Executable not found. Please build the project first.
        exit /b 1
    )
)

if "%1" == "migrate" (
    if "%2" == "" (
        echo Please specify a name for the migration.
        exit /b 1
    )
    echo Creating new migration %2...
    %MIGRATE_CMD% create -ext sql -dir %MIGRATION_DIR% %2
)

if "%1" == "migrate-up" (
    echo Running Go up migration...
    go run cmd\migrate\main.go up
)

if "%1" == "migrate-down" (
    echo Running Go down migration...
    go run cmd\migrate\main.go down
)

endlocal
