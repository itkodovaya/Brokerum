#!/bin/bash

# Запуск TenderHelp проекта

echo "Запуск TenderHelp..."

# Проверка наличия Go
if ! command -v go &> /dev/null; then
    echo "Ошибка: Go не установлен. Установите Go версии 1.21+"
    exit 1
fi

# Проверка наличия Node.js
if ! command -v node &> /dev/null; then
    echo "Ошибка: Node.js не установлен. Установите Node.js"
    exit 1
fi

# Установка зависимостей Go
echo "Установка зависимостей Go..."
go mod tidy

# Установка зависимостей React
echo "Установка зависимостей React..."
cd frontend
npm install
cd ..

# Запуск Go сервера в фоне
echo "Запуск Go сервера..."
go run main.go &
GO_PID=$!

# Ожидание запуска сервера
sleep 3

# Запуск React приложения
echo "Запуск React приложения..."
cd frontend
npm start &
REACT_PID=$!

echo "Проект запущен!"
echo "Go сервер: http://localhost:8080"
echo "React приложение: http://localhost:3000"
echo ""
echo "Для остановки нажмите Ctrl+C"

# Ожидание сигнала завершения
trap "kill $GO_PID $REACT_PID; exit" INT TERM
wait
