# Calculation API Service

REST API сервис для вычисления математических выражений, написанный на Go с использованием фреймворка Echo и ORM GORM (PostgreSQL).

## 🚀 Возможности
- Создание вычислений (POST /calculations)
- Получение списка всех вычислений (GET /calculations)
- Обновление выражения (PATCH /calculations/:id)
- Удаление вычисления (DELETE /calculations/:id)
- Валидация входных данных (формат UUID, пустые выражения)
- Корректная обработка ошибок (400 Bad Request, 404 Not Found, 500 Internal Server Error)

## 🛠 Стек технологий
- **Язык:** Go (Golang)
- **Веб-фреймворк:** Echo v4
- **База данных:** PostgreSQL
- **ORM:** GORM
- **Парсинг выражений:** govaluate

## 📦 Как запустить локально

1. Убедитесь, что у вас установлены Go и PostgreSQL.
2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/grammbatanov-maker/awesomeProject1.git
   cd awesomeProject1

Настройте подключение к базе данных в файле internal/db/db.go (измените dsn под свои credentials).

Запустите сервер:
go run cmd/main.go

Сервер будет доступен по адресу: http://localhost:8080

📡 Примеры запросов
Создать вычисление:

bash
curl -X POST http://localhost:8080/calculations \
  -H "Content-Type: application/json" \
  -d '{"expression": "2 + 2 * 2"}'
Получить все вычисления:
curl http://localhost:8080/calculations
