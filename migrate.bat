@echo off
setlocal enabledelayedexpansion

REM Pindah ke direktori skrip ini berada (auth/)
cd /d "%~dp0"

REM Baca file .env dari folder saat ini
for /f "tokens=1* delims==" %%A in ('findstr /v "^#" .\.env') do (
    set "%%A=%%B"
)

REM Buat connection string DB
set "DB_URL=postgres://%DB_USER%:%DB_PASSWORD%@%DB_HOST%:%DB_PORT%/%DB_NAME%?sslmode=disable"

if "%~1"=="up" goto do_up
if "%~1"=="down" goto do_down
if "%~1"=="force" goto do_force
if "%~1"=="create" goto do_create
goto do_help

:do_up
echo Menjalankan migrate up...
migrate.exe -path db/migrations -database "%DB_URL%" up
goto :eof

:do_down
echo Menjalankan migrate down...
migrate.exe -path db/migrations -database "%DB_URL%" down
goto :eof

:do_force
if "%~2"=="" (
    echo Error: Butuh parameter versi. Contoh: migrate.bat force 1
    exit /b 1
)
echo Memaksa versi migrasi ke %~2...
migrate.exe -path db/migrations -database "%DB_URL%" force %~2
goto :eof

:do_create
if "%~2"=="" (
    echo Error: Butuh nama migrasi. Contoh: migrate.bat create init_schema
    exit /b 1
)
echo Membuat file migrasi baru: %~2...
migrate.exe create -ext sql -dir db/migrations -seq %~2
goto :eof

:do_help
echo ===== GOLANG MIGRATE HELPER =====
echo Penggunaan:
echo migrate.bat up                 - Jalankan semua migrasi yang belum dijalankan (Up)
echo migrate.bat down               - Revert/rollback semua migrasi (Down)
echo migrate.bat force [version]    - Paksa versi skema migrasi jika error (contoh: migrate.bat force 1)
echo migrate.bat create [name]      - Buat file migrasi baru (contoh: migrate.bat create add_users_table)
echo.
echo Pastikan Anda sudah menginstall golang-migrate CLI
goto :eof
