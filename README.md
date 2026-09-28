# AuthService

Сервис авторизации для магазина бытовой техники.
Отдельное приложение со своей базой данных PostgreSQL, которое взаимодействует с внешним миром **только по gRPC** (HTTP/2).

Основное приложение магазина: [ShopAPI](https://github.com/sonofsun61/ShopAPI). Оно проксирует в AuthService запросы `/register`, `/auth`, `/reset`, а на остальных ручках проверяет JWT-токен через middleware.

## Возможности

| Метод gRPC       | Что делает                                                        |
|------------------|-------------------------------------------------------------------|
| `Register`       | Создаёт пользователя (почта, имя, фамилия, телефон, пароль), возвращает подписанный токен |
| `Login`          | Проверяет пару логин/пароль, возвращает токен                     |
| `ChangePassword` | Меняет пароль пользователя (нужен старый пароль)                  |
| `ResetPassword`  | Восстановление пароля: новый пароль «отправляется на почту» (заглушка, вывод в консоль сервиса) |

## Стек

- Go, gRPC (`google.golang.org/grpc`), Protocol Buffers
- PostgreSQL 16 (в Docker), драйвер `pgx/v5`
- Миграции: `golang-migrate`

## Структура проекта

```
.
├── api/proto/auth/v1/auth.proto   # контракт gRPC
├── cmd/authservice/main.go        # точка входа
├── gen/authpb/                    # сгенерированный код (protoc)
├── internal/grpcserver/           # transport-слой: реализация gRPC-методов
├── migrations/                    # SQL-миграции
├── docker-compose.yml             # PostgreSQL для сервиса
└── .env                           # секреты (не коммитится)
```

## Быстрый старт

### 1. Требования

- Go
- Docker и Docker Compose
- `protoc` с плагинами `protoc-gen-go` и `protoc-gen-go-grpc`
- CLI `migrate` (`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)
- (необязательно) `grpcurl` для ручной проверки

### 2. Настройка окружения

Создайте файл `.env` в корне проекта (он должен быть в `.gitignore`):

```env
POSTGRES_USER=<пользователь>
POSTGRES_PASSWORD=<пароль>
POSTGRES_DB=auth_db
DATABASE_URL=postgres://<пользователь>:<пароль>@localhost:5433/auth_db?sslmode=disable
```

Порт хоста `5433` выбран, чтобы не конфликтовать с базой ShopAPI. Загрузить переменные в терминал:

```bash
set -a; source .env; set +a
```

### 3. База данных

```bash
docker compose up -d
migrate -path migrations -database "$DATABASE_URL" up
```

Откат последней миграции: `migrate -path migrations -database "$DATABASE_URL" down 1`.

### 4. Генерация кода из `.proto`

```bash
protoc \
  --proto_path=api/proto \
  --go_out=. --go_opt=module=authservice \
  --go-grpc_out=. --go-grpc_opt=module=authservice \
  api/proto/auth/v1/auth.proto
```

### 5. Запуск

```bash
go run ./cmd/authservice
```

Сервис слушает порт `50051`.
