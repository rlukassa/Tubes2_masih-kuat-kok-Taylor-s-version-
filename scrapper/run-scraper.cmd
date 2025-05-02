@echo off

echo [*] Menghapus bekas kompilasi sebelumnya...
del /f /q *.csv
del /f /q *.db
del /f /q go.mod
del /f /q go.sum

echo [*] Memeriksa dan mengunduh dependensi Go...
if not exist go.mod (
    go mod init scrape_elements
)
go mod tidy

echo [*] Menjalankan scraper...
go run ./scraper.go

echo.
pause
