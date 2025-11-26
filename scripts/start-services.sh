#!/bin/bash

# Скрипт для запуска всех сервисов локально

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Функция для запуска сервиса
start_service() {
  local service=$1
  local service_dir="$ROOT_DIR/$service"
  local log_file="/tmp/${service}.log"
  local pid_file="/tmp/${service}.pid"
  
  # Проверяем и останавливаем старый процесс, если он существует
  if [ -f "$pid_file" ]; then
    local old_pid=$(cat "$pid_file")
    if ps -p "$old_pid" > /dev/null 2>&1; then
      echo "🛑 Останавливаю старый процесс $service (PID: $old_pid)..."
      kill "$old_pid" 2>/dev/null || kill -9 "$old_pid" 2>/dev/null || true
      sleep 1
    fi
    rm -f "$pid_file"
  fi
  
  # Проверяем, не занят ли порт другим процессом
  local port=""
  case $service in
    inventory)
      port="50051"
      ;;
    payment)
      port="50052"
      ;;
    order)
      port="8080"
      ;;
    notification)
      port="8081"
      ;;
    iam)
      port="50056"
      ;;
  esac
  
  if [ -n "$port" ]; then
    local port_pid=$(lsof -ti :$port 2>/dev/null | head -1)
    if [ -n "$port_pid" ]; then
      echo "⚠️  Порт $port занят процессом $port_pid. Останавливаю..."
      kill "$port_pid" 2>/dev/null || kill -9 "$port_pid" 2>/dev/null || true
      sleep 1
      # Проверяем еще раз и принудительно завершаем, если нужно
      local still_running=$(lsof -ti :$port 2>/dev/null | head -1)
      if [ -n "$still_running" ]; then
        echo "⚠️  Принудительно завершаю процесс на порту $port..."
        kill -9 "$still_running" 2>/dev/null || true
        sleep 1
      fi
    fi
  fi
  
  echo "🚀 Запускаю $service..."
  
  # Очищаем старый лог перед запуском
  > "$log_file"
  
  # Запускаем сервис с перенаправлением вывода
  # Важно: запускаем из корня проекта, чтобы пути к .env файлам были правильными
  # Используем nohup для правильной работы в фоне
  # Явно указываем рабочую директорию через -C опцию или через cd в подпроцессе
  cd "$ROOT_DIR"
  nohup go run "$service_dir/cmd/main.go" >> "$log_file" 2>&1 &
  local pid=$!
  echo $pid > "$pid_file"
  
  # Даем процессу немного времени на запуск и запись ошибок в лог
  sleep 2
  
  # Проверяем, что процесс еще работает
  if ! ps -p "$pid" > /dev/null 2>&1; then
    echo "❌ $service завершился сразу после запуска. Проверьте лог: $log_file"
    # Ждем еще немного, чтобы убедиться, что все данные записались в лог
    sleep 1
    if [ -s "$log_file" ]; then
      echo "Последние строки лога:"
      tail -30 "$log_file"
    else
      echo "⚠️  Лог файл пуст - возможно, ошибка при загрузке конфигурации или инициализации"
      echo "💡 Попробуйте запустить сервис вручную: cd $ROOT_DIR && go run $service_dir/cmd/main.go"
    fi
    rm -f "$pid_file"
    return 1
  fi
  
  # Ждем еще немного, чтобы процесс успел полностью запуститься
  # IAM требует больше времени для подключения к БД и выполнения миграций
  if [ "$service" = "iam" ]; then
    sleep 6
  elif [ "$service" = "order" ]; then
    # Order также требует времени для подключения к БД
    sleep 3
  else
    sleep 1
  fi
  
  # Проверяем, что порт слушается (более надежный способ проверки)
  if [ -n "$port" ]; then
    local listening_pid=$(lsof -ti :$port 2>/dev/null | head -1)
    if [ -n "$listening_pid" ]; then
      echo "✅ $service запущен (PID процесса на порту $port: $listening_pid, лог: $log_file)"
      echo $listening_pid > "$pid_file"
      return 0
    else
      echo "❌ $service не запустился (порт $port не слушается). Проверьте лог: $log_file"
      tail -30 "$log_file"
      rm -f "$pid_file"
      return 1
    fi
  else
    # Если порт не определен, проверяем процесс
    if ps -p "$pid" > /dev/null 2>&1; then
      echo "✅ $service запущен (PID: $pid, лог: $log_file)"
      return 0
    else
      echo "❌ $service завершился с ошибкой. Проверьте лог: $log_file"
      tail -30 "$log_file"
      rm -f "$pid_file"
      return 1
    fi
  fi
}

# Проверяем, что .env файлы существуют
if [ ! -f "$ROOT_DIR/deploy/compose/inventory/.env" ] || [ ! -f "$ROOT_DIR/deploy/compose/order/.env" ] || [ ! -f "$ROOT_DIR/deploy/compose/payment/.env" ] || [ ! -f "$ROOT_DIR/deploy/compose/assembly/.env" ] || [ ! -f "$ROOT_DIR/deploy/compose/notification/.env" ] || [ ! -f "$ROOT_DIR/deploy/compose/iam/.env" ]; then
  echo "📝 .env файлы не найдены. Генерирую..."
  cd "$ROOT_DIR"
  if command -v task &> /dev/null; then
    task env:generate || echo "⚠️  Не удалось сгенерировать .env файлы через task. Убедитесь, что они существуют в deploy/compose/{service}/.env"
  else
    echo "⚠️  task не найден. Убедитесь, что .env файлы существуют в deploy/compose/{service}/.env"
  fi
fi

# Запускаем сервисы в правильном порядке
# Используем переменную для отслеживания ошибок, но не останавливаем скрипт при ошибке одного сервиса
ERRORS=0

start_service inventory || ERRORS=$((ERRORS + 1))
sleep 1
start_service payment || ERRORS=$((ERRORS + 1))
sleep 1
start_service order || ERRORS=$((ERRORS + 1))
sleep 1
start_service assembly || ERRORS=$((ERRORS + 1))
sleep 1
start_service notification || ERRORS=$((ERRORS + 1))
sleep 1
start_service iam || ERRORS=$((ERRORS + 1))

echo
if [ $ERRORS -eq 0 ]; then
  echo "🎉 Все сервисы запущены!"
else
  echo "⚠️  Некоторые сервисы не запустились (ошибок: $ERRORS)"
  echo "💡 Проверьте логи в /tmp/*.log для деталей"
fi
echo "Для остановки используйте: ./scripts/stop-services.sh"

