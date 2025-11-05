#!/bin/bash

# Скрипт для остановки всех сервисов

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

stop_service() {
  local service=$1
  local pid_file="/tmp/${service}.pid"
  
  if [ ! -f "$pid_file" ]; then
    echo "⚠️  Сервис $service не запущен (файл PID не найден)"
    return 0
  fi
  
  local pid=$(cat "$pid_file")
  if ! ps -p "$pid" > /dev/null 2>&1; then
    echo "⚠️  Процесс $service (PID: $pid) не найден"
    rm -f "$pid_file"
    return 0
  fi
  
  echo "🛑 Останавливаю $service (PID: $pid)..."
  kill "$pid" 2>/dev/null || true
  
  # Ждем завершения процесса
  for i in {1..10}; do
    if ! ps -p "$pid" > /dev/null 2>&1; then
      break
    fi
    sleep 0.5
  done
  
  # Если процесс все еще работает, принудительно завершаем
  if ps -p "$pid" > /dev/null 2>&1; then
    echo "⚠️  Принудительно завершаю $service..."
    kill -9 "$pid" 2>/dev/null || true
  fi
  
  rm -f "$pid_file"
  echo "✅ $service остановлен"
}

# Останавливаем сервисы в обратном порядке
stop_service order
stop_service payment
stop_service inventory

echo
echo "🎉 Все сервисы остановлены!"




