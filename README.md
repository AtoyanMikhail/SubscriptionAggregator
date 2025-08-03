# Subscription Aggregator

REST-сервис для агрегации данных об онлайн-подписках пользователей.

## Описание

Сервис предоставляет API для управления подписками пользователей с возможностью создания, чтения, обновления, удаления и листинга подписок, а также расчета общей стоимости подписок за выбранный период.

## Архитектура

Проект следует принципам чистой архитектуры с разделением на слои:

```
├── cmd/                    # Точка входа в приложение
├── internal/              # Внутренняя логика приложения
│   ├── config/           # Конфигурация
│   ├── handlers/         # HTTP обработчики
│   ├── logger/           # Логирование
│   ├── repository/       # Слой доступа к данным
│   └── service/          # Бизнес-логика
├── migrations/           # Миграции базы данных
├── configs/             # Конфигурационные файлы
├── docs/               # Swagger документация
└── docker-compose.local.yml
```

## Требования

- Go 1.24.4+
- Docker & Docker Compose
- PostgreSQL 15+

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/AtoyanMikhail/SubscriptionAggregator.git
cd SubscriptionAggregator
```

### 2. Запуск с помощью Docker Compose

```bash
# Запуск всех сервисов
docker-compose -f docker-compose.local.yml up -d

# Проверка статуса
docker-compose -f docker-compose.local.yml ps
```

### 3. Альтернативный запуск для разработки

```bash
# Установка зависимостей
go mod tidy

# Генерация Swagger документации
make swagger

# Запуск PostgreSQL
make docker-up

# Запуск приложения
make run
```

## Конфигурация

Конфигурация приложения находится в файле `configs/app/config_local.yaml`:

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  db_name: "subscription_aggregator"
  ssl_mode: "disable"
```

## API Endpoints

### Базовый URL
```
http://localhost:8080
```

### Swagger документация
```
http://localhost:8080/swagger/index.html
```

### Endpoints

#### 1. CRUD операции с подписками

**Создание подписки**
```http
POST /api/v1/subscriptions
Content-Type: application/json

{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

**Получение подписки**
```http
GET /api/v1/subscriptions/{user_id}/{subscription_id}
```

**Обновление подписки**
```http
PUT /api/v1/subscriptions/{user_id}/{subscription_id}
Content-Type: application/json

{
  "service_name": "Netflix Premium",
  "price": 599,
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

**Удаление подписки**
```http
DELETE /api/v1/subscriptions/{user_id}/{subscription_id}
```

**Получение всех подписок пользователя**
```http
GET /api/v1/subscriptions/user/{user_id}
```

#### 2. Расчет стоимости подписок

**Расчет общей стоимости за период**
```http
GET /api/v1/subscriptions/cost?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&start_date=01-2025&end_date=12-2025&service_names=Netflix&service_names=Spotify
```

**Параметры запроса:**
- `user_id` (обязательный) - UUID пользователя
- `start_date` (обязательный) - начало периода в формате MM-YYYY
- `end_date` (обязательный) - конец периода в формате MM-YYYY
- `service_names` (опциональный) - массив названий сервисов для фильтрации

#### 3. Health Check

```http
GET /health
```

## Модель данных

### Подписка (Subscription)

| Поле         | Тип     | Описание                              |
|--------------|---------|---------------------------------------|
| id           | SERIAL  | Уникальный идентификатор              |
| service_name | TEXT    | Название сервиса                      |
| price        | INTEGER | Стоимость в рублях                    |
| user_id      | TEXT    | UUID пользователя                     |
| start_date   | DATE    | Дата начала подписки                  |
| end_date     | DATE    | Дата окончания подписки (опционально) |

### Индексы

- `idx_subscriptions_user_id` - для быстрого поиска по пользователю
- `idx_subscriptions_service_name` - для фильтрации по сервису
- `idx_subscriptions_start_date` - для поиска по дате начала
- `idx_subscriptions_end_date` - для поиска по дате окончания

## Особенности реализации

### Валидация данных
- UUID формат для user_id
- Неотрицательные значения для price
- Формат дат MM-YYYY для start_date и end_date
- Обязательные поля: service_name, price, user_id, start_date

### Расчет стоимости
- Учитываются только полные месяцы подписки
- Если end_date не указана, подписка считается бессрочной
- Поддерживается фильтрация по конкретным сервисам
- Возвращается детальная разбивка по каждой подписке

## Разработка

### Установка инструментов разработки

```bash
make install-tools
```

### Команды для разработки

```bash
# Форматирование кода
make fmt

# Линтинг
make lint

# Запуск тестов
make test

# Генерация Swagger документации
make swagger

# Полная настройка окружения
make setup
```

## Docker

### Сборка образа

```bash
docker build -t subscription-aggregator .
```

### Запуск контейнеров

```bash
# Разработка
docker-compose -f docker-compose.yml up -d
```

## Мониторинг

### Health Check
```http
GET /health
```

### Логи
Логи приложения по умолчанию доступны через Docker:
```bash
docker-compose -f docker-compose.local.yml logs -f app
```
