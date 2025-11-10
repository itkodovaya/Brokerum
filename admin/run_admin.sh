#!/bin/bash
echo "Запуск админки TenderHelp..."
cd "$(dirname "$0")"
go run main.go
