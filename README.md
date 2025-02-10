# Магазин мерча

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