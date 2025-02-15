# Магазин мерча


[![Merch App CI](https://www.github.com/justcgh9/merch_store/actions/workflows/ci.yml/badge.svg)](https://www.github.com/justcgh9/merch_store/actions/workflows/ci.yml) [![Coverage Status](https://coveralls.io/repos/github/justcgh9/merch_store/badge.svg)](https://coveralls.io/github/justcgh9/merch_store)
## Структура проекта

**Ниже** Вы можете видеть структуру проекта:

---

```bash
    merch_store/
    ├── cmd
    │   └── merch-store
    │       └── main.go
    ├── config
    │   └── local.yml
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    └── internal
        ├── config
        │   └── config.go
        ├── http-server
        │   ├── handlers
        │   │   ├── auth
        │   │   │   └── auth.go
        │   │   ├── buy
        │   │   │   └── buy.go
        │   │   ├── info
        │   │   │   └── info.go
        │   │   └── send
        │   │       └── send.go
        │   └── middleware
        │       └── auth
        │           └── auth.go
        ├── models
        │   ├── inventory
        │   │   └── inventory.go
        │   ├── transaction
        │   │   └── transaction.go
        │   └── user
        │       └── user.go
        ├── services
        │   ├── coin
        │   │   └── coin.go
        │   ├── merch
        │   │   └── merch.go
        │   └── user
        │       └── user.go
        └── storage
            ├── postgres
            │   └── postgres.go
            └── storage.go

```

## Настройка базы данных и миграций

### Схема базы данных

В нашем PostgreSQL-хранилище используется четыре основные таблицы:

- **Users** – хранит имена пользователей и их пароли.
- **Balance** – отслеживает баланс пользователей.
- **Inventory** – содержит 10 различных счетчиков товаров для каждого пользователя.
- **History** – фиксирует транзакции между пользователями.

Простая схема базы данных:

```
    +-------------+       +------------+ 
    |   Users     |       |  Balance   |
    |-------------|       |------------|
--> | username PK | <---> | username PK|
|   | password    |       | balance    |
|   +-------------+       +------------+
|
|   +------------+       +---------------------------+
|   | Inventory  |       |         History           |
|   |------------|       |---------------------------|
--> | username PK| <---> | from_user FK -> Users     |
    | t-shirt    |       | to_user   FK -> Users     |
    | cup        |       | amount                    |
    | book       |       +---------------------------+
    | ...        |
    +------------+
```

### Миграции

Для управления миграциями используется **golang-migrate**.

#### Применение миграций

Чтобы применить миграции, выполните команду:

```sh
go run cmd/migrator/main.go -db "postgres://merch_user:merch_password@localhost:5432/merch_db?sslmode=disable" -path "./migrations" -action up
```

#### Откат миграций

Для отката миграций используйте:

```sh
go run cmd/migrator/main.go -db "postgres://merch_user:merch_password@localhost:5432/merch_db?sslmode=disable" -path "./migrations" -action down
```

### Запуск PostgreSQL через Docker

Запустить локальную базу данных PostgreSQL можно с помощью команды:

```sh
docker-compose up -d postgres
```

Этот запуск использует конфигурацию из `docker-compose.yml`. Убедитесь, что порт 5432 свободен, или замените его на другой в ямле.

---

С этим набором инструментов вы можете легко управлять схемой базы данных и применять изменения при необходимости. 🚀


