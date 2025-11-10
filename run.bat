@echo off
echo Запуск Brokerum...

REM Проверка наличия Go
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo Ошибка: Go не установлен. Установите Go версии 1.21+
    pause
    exit /b 1
)

REM Проверка наличия Node.js
where node >nul 2>nul
if %errorlevel% neq 0 (
    echo Ошибка: Node.js не установлен. Установите Node.js
    pause
    exit /b 1
)

REM Установка зависимостей Go
echo Установка зависимостей Go...
go mod tidy

REM Установка зависимостей React
echo Установка зависимостей React...
cd frontend
call npm install
cd ..

REM Запуск Go сервера в фоне
echo Запуск Go сервера...
start /b go run main.go

REM Ожидание запуска сервера
timeout /t 3 /nobreak >nul

REM Запуск React приложения
echo Запуск React приложения...
cd frontend
start /b npm start
cd ..

echo Проект запущен!
echo Go сервер: http://localhost:8080
echo React приложение: http://localhost:3000
echo.
echo Для остановки закройте это окно
pause
