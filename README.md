# **Тестовое задание для стажёра Backend (осенняя волна 2025)**

## **Сервис назначения ревьюеров для Pull Request'ов**

### Описание проекта

Микросервис для автоматизации процесса назначения ревьюеров на Pull Request'ы. Сервис управляет командами разработчиков, отслеживает активность пользователей и интеллектуально распределяет нагрузку по code review между членами команды.

**Основной функционал:**
- Управление командами и участниками
- Автоматическое назначение до 2 ревьюеров из команды автора PR
- Переназначение ревьюеров с учётом текущей нагрузки
- Статистика по назначениям и Pull Request'ам
- Отслеживание активности пользователей
- Просмотр PR, назначенных конкретному ревьюеру

### Архитектура

Проект следует принципам Clean Architecture с чёткой организацией слоёв:

```
.
├── cmd/app/  # ------------------------- Точка входа приложения
├── internal/
│   ├── config/  # ---------------------- Конфигурация приложения
│   ├── domain/  # ---------------------- Доменные модели (Team, User, PullRequest)
│   ├── dto/  # ------------------------- Data Transfer Objects для API
│   ├── handler/  # --------------------- HTTP обработчики (transport layer)
│   ├── service/  # --------------------- Бизнес-логика
│   │   ├── pull_request_service.go
│   │   ├── team_service.go
│   │   ├── user_service.go
│   │   └── reviewer.go  # -------------- Алгоритм выбора ревьюеров
│   ├── repository/  # ------------------ Интерфейсы репозиториев
│   │   └── postgres/ # ----------------- Реализация для PostgreSQL
│   └── transport/ # -------------------- HTTP роутер и middleware
├── migrations/ # ----------------------- SQL миграции базы данных
├── docs/  # ---------------------------- OpenAPI спецификация и Swagger
└── pkg/logger/ # ----------------------- Логирование
```

**Технологический стек:**
- **Go 1.24**
- **Golangci-lint** - линтер
- **Gorilla Mux** - HTTP роутер
- **PostgreSQL** - база данных
- **pgx** - PostgreSQL драйвер
- **testify/suite** - тестирование
- **Swagger** - документация API
- **Docker & Docker Compose** - контейнеризация

### Быстрый старт
Перед запуском необходимо заполнить `.env` файл по шаблону ([`.env.example`](.env.example)).

#### Локальная сборка
```bash
# Сборка бинарника
go build -v ./cmd/app/main.go # или make build

# Запуск
./main
# Сервер запущен на порте :8080
```

#### Запуск через Docker Compose

```bash
# Сборка образов
docker compose -f deploy/docker-compose.yml build

# Запуск сервисов (приложение + PostgreSQL)
docker compose --env-file deploy/.env -f deploy/docker-compose.yml up

# Запуск в фоновом режиме
docker compose --env-file deploy/.env -f deploy/docker-compose.yml up -d # или make up

# Остановка
docker compose -f deploy/docker-compose.yml down  # или make down

# Прогон миграций
make migrate-up
```

После запуска сервис доступен по адресу `http://localhost:8080`

#### Makefile

```bash
# Сборка
make build                    # Компиляция бинарника ./main
make up                       # Запуск проекта в Docker контейнере и прогон миграций
make down                     # Остановка Docker контейнера

# Тестирование
make test                     # Запуск всех тестов (юнит, нагрузочные)

make test-unit                # Запуск юнит-тестов
make test-unit-coverage       # Тесты с покрытием (консольный вывод)
make test-unit-coverage-html  # Тесты с HTML отчётом покрытия

# Нагрузочное тестирование
make test-load                # Все нагрузочные тесты
make test-load-team           # Нагрузка на эндпоинты команд
make test-load-pr             # Нагрузка на эндпоинты PR
make test-load-all            # Комплексный нагрузочный тест


# Миграции базы данных
make migrate-up               # Применить все миграции
make migrate-down             # Откатить последнюю миграцию
make migrate-down-all         # Откатить все миграции
make migrate-version          # Показать текущую версию
make migrate-create           # Создать новую миграцию (интерактивно)
make migrate-force            # Принудительно установить версию

# Миграции в Docker
make migrate-docker-up        # Применить миграции в Docker контейнере
make migrate-docker-down      # Откатить миграцию в Docker

# База данных
make db-shell                 # Открыть psql в Docker контейнере

# Линтер
make lint                     # Линтер (golangci-lint)
make lint-fix                 # Запуск линтера с флагом --fix

# Очистка
make clean                    # Удалить сгенерированные файлы (*.out, *.html, main)
```

### Документация API

**Swagger** доступен по адресу: `http://localhost:8080/swagger/index.html`

#### Доступные эндпоинты

**Управление командами:**
- `POST /team/add` - создание команды с участниками
  ```json
  {
    "team_name": "backend",
    "members": ["user1", "user2", "user3"]
  }
  ```
- `GET /team/get?team_name=backend` - получение информации о команде

**Управление пользователями:**
- `POST /users/setIsActive` - изменение статуса активности пользователя
  ```json
  {
    "user_id": "user1",
    "is_active": false
  }
  ```
- `GET /users/getReview?user_id=user1` - список PR, где пользователь назначен ревьюером

**Управление Pull Request'ами:**
- `POST /pullRequest/create` - создание PR (автоматически назначает до 2 ревьюеров)
  ```json
  {
    "pull_request_id": "pr123",
    "pull_request_name": "Add new feature",
    "author_id": "user1"
  }
  ```
- `POST /pullRequest/merge` - отметка PR как смёрженного (идемпотентная операция)
  ```json
  {
    "pull_request_id": "pr123"
  }
  ```
- `POST /pullRequest/reassign` - переназначение ревьюера
  ```json
  {
    "pull_request_id": "pr123",
    "old_user_id": "user2"
  }
  ```
- `GET /pullRequest/get?pull_request_id=pr123` - получение информации о PR

**Статистика:**
- `GET /stats` - статистика по назначениям (сколько раз каждый пользователь был назначен, сколько ревьюеров у каждого PR)

**Системные эндпоинты:**
- `GET /health` - проверка здоровья сервиса
- `GET /` - корневая страница

#### Регенерация Swagger документации

После изменения API аннотаций в коде:
```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/app/main.go -o docs
```

### Тестирование

Проект покрыт юнит-тестами и интеграционными тестами с использованием **testify/suite** и table-driven подхода.

#### Запуск тестов

```bash
# Все юнит-тесты (internal/...)
make test-unit

# С подробным выводом
go test -v ./internal/...

# Только тесты сервисного слоя
go test -v ./internal/service/...

# Только тесты handler'ов
go test -v ./internal/handler/...

# Запуск с race detector
go test -race ./internal/...
```

#### Покрытие кода тестами

```bash
# Просмотр покрытия в консоли
make test-unit-coverage

# Генерация HTML отчёта
make test-unit-coverage-html
# Файл coverage.html будет создан автоматически
open coverage.html
```

**Текущее покрытие:**
- **Сервисный слой:** 87.3% (4 test suite, 44+ тестовых случая)
- **Handler слой:** 79.5% (5 test suite, 22+ тестовых случая)

#### Структура тестов

**Service тесты** (`internal/service/*_test.go`):
- Юнит-тесты бизнес-логики
- Используют mock-репозитории для изоляции
- Организованы в testify suites:
  - `PullRequestServiceTestSuite` - тесты создания, мёрджа, переназначения PR
  - `TeamServiceTestSuite` - тесты создания и получения команд
  - `UserServiceTestSuite` - тесты управления пользователями
  - `ReviewerTestSuite` - тесты алгоритмов выбора ревьюеров

**Handler тесты** (`internal/handler/*_test.go`):
- Интеграционные HTTP тесты
- Используют реальные сервисы с mock-репозиториями
- Организованы в testify suites:
  - `PRHandlerTestSuite` - тесты HTTP эндпоинтов PR
  - `TeamHandlerTestSuite` - тесты эндпоинтов команд
  - `UserHandlerTestSuite` - тесты эндпоинтов пользователей
  - `HealthHandlerTestSuite` - тесты системных эндпоинтов
  - `StatsHandlerTestSuite` - тесты эндпоинта статистики

**Примеры запуска:**
```bash
# Запуск конкретного test suite
go test -v ./internal/service -run TestPullRequestService

# Запуск конкретного тестового метода
go test -v ./internal/service -run TestPullRequestService/TestCreate

# Все тесты конкретного пакета
go test -v ./internal/handler
```

### Нагрузочное тестирование

Для нагрузочного тестирования используется **Vegeta**.

#### Быстрый старт

```bash
# Установка Vegeta (macOS)
brew install vegeta

# Запуск сервиса
docker compose --env-file .env up -d

# Запуск всех нагрузочных тестов
make test-load

# Отдельные наборы тестов
make test-load-team    # Тесты эндпоинтов команд
make test-load-pr      # Тесты эндпоинтов PR
make test-load-all     # Комплексный тест всех эндпоинтов
```

**Целевые показатели производительности:**
- p95 латентность < 100ms
- p99 латентность < 200ms
- Успешность запросов > 99%
- RPS: 100-1000 в зависимости от эндпоинта

*Отчет по результатам тестирования расположен в [loadtest/report.md](loadtest/report.md).*

### База данных

Проект использует **PostgreSQL** с миграциями для управления схемой.

#### Структура БД

**Таблицы:**
- `teams` - команды разработчиков
- `users` - пользователи с привязкой к командам и флагом активности
- `pull_requests` - Pull Request'ы с назначенными ревьюерами

#### Миграции

Миграции находятся в директории `migrations/`:
```
migrations/
├── 000001_teams.up.sql                 # Создание таблицы teams
├── 000001_teams.down.sql               # Откат миграции teams
├── 000002_users.up.sql                 # Создание таблицы users
├── 000002_users.down.sql               # Откат миграции users
├── 000003_pull_requests.up.sql         # Создание таблицы pull_requests
└── 000003_pull_requests.down.sql       # Откат миграции pull_requests
```

#### Подключение к БД

Параметры подключения задаются через переменные окружения:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=avito_test
```

### Бизнес-логика

#### Алгоритм назначения ревьюеров

При создании PR сервис автоматически назначает до 2 ревьюеров:

1. **Фильтрация кандидатов:**
   - Только активные пользователи (`is_active = true`)
   - Только члены команды автора
   - Исключается сам автор

2. **Случайный выбор:**
   - Выбирается до 2 случайных кандидатов
   - Если в команде меньше 2 подходящих кандидатов, назначаются все доступные

3. **Обработка ошибок:**
   - Если нет активных кандидатов → ошибка `NO_CANDIDATE`
   - Если автор не найден → ошибка `NOT_FOUND`
   - Если команда не существует → ошибка `NOT_FOUND`

#### Переназначение ревьюера

При переназначении (`/pullRequest/reassign`):

1. **Валидация:**
   - PR не должен быть смёрджен
   - Старый ревьюер должен быть назначен на этот PR

2. **Выбор замены:**
   - Ищутся активные члены команды старого ревьюера
   - Исключаются уже назначенные ревьюеры
   - Исключается автор PR
   - Выбирается случайный кандидат

3. **Атомарное обновление:**
   - Старый ревьюер удаляется из списка
   - Новый ревьюер добавляется
   - PR обновляется в БД

### Мои комментарии 
В процессе разработки были допущены некоторые упущения, связанные с близким дедлайном и недостатком опыта.
К ним относятся:
- отсутствие транзакций при работе с БД
- возможные ошибки в архитектуре приложения. Например, миграции при запуске в контейнере в `docker-compose.yml` запихнуть так и не получилось)
- не самый оптимальный способ сбора статистики - каждый раз обращаемся к БД и считаем по новой...

и т.д. 

В любом случае, удовольствие в процессе разработки решения было получено, результатом я доволен) 😄

---

Разработано в рамках тестового задания для позиции Cтажёр-Backend в Авито.
