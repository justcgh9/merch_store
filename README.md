# Магазин мерча


[![Merch App CI](https://www.github.com/justcgh9/merch_store/actions/workflows/ci.yml/badge.svg)](https://www.github.com/justcgh9/merch_store/actions/workflows/ci.yml) [![Coverage Status](https://coveralls.io/repos/github/justcgh9/merch_store/badge.svg)](https://coveralls.io/github/justcgh9/merch_store)
## Структура проекта

Проект написан на golang + postgresql. Я решил условно разделить домен на три части: пользователь, инвентарь (всё, что касается мерча), и деньги. Также стоит упомянуть о следующих вещах:

1. В качестве роутера используется `chi`.
1. Для работы с базой данных используется `pgx`.
1. Для удобной работы настроены миграции.
1. Проект запускается через `Docker` || `Docker-compose`. 
1. Информация об окружении берется из конфига, например [local](/config/local.yml).
1. Необходимые переменные окружения можно найти в [docker-compose.yml](/docker-compose.yml), но для них установлены значения по умолчанию в [Dockerfile](/Dockerfile)
1. По умолчанию, docker-compose будет запускать два контейнера - один из которых с бд, поэтому необходимо правильно настроить порты и строку подключение в конфиге (если 5432 порт занят)
1. Все необходимые команды включены в [таскфайл](/Taskfile.yml)

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
task migrate-up
```

#### Откат миграций

Для отката миграций используйте:

```sh
task migrate-down
```

### Запуск PostgreSQL через Docker

Запустить локальную базу данных PostgreSQL можно с помощью команды:

```sh
docker-compose up -d postgres
```

Этот запуск использует конфигурацию из `docker-compose.yml`. Убедитесь, что порт 5432 свободен, или замените его на другой в ямле.

---

С этим набором инструментов вы можете легко управлять схемой базы данных и применять изменения при необходимости. 🚀

## Нагрузочное тестирование

### Объяснение

Я провел небольшое нагрузочное тестирование с помощью **Locust**. В папку с миграциями я поместил `.sql` файлы которые использовал для создания данных фейковых пользователей нагрузочного тестирования. В `stress_tesing/locustfile.py` я описал необходимые действия для каждой из тасок. Далее с помощь команды 

```sh
locust
```

Я запустил тестер, перешел по 8089 порту на локалхосте, и там с помощью графического интерфейса настроил необходимые данные о нагрузке. Стоит заметить, что всем пользователям я выдал по много денег на баланс, что должно было убрать ошибочные http коды от недостаточных средств, но это могло только повысить нагрузку (незначительно), так как успешный сценарий длиннее ошибочного. Помимо этого, хочу отметить две вещи:

1. Локуст съедал бОльшую часть CPU, моментами оставляя моё приложение голодать
2. Несмотря на то, что я постарался закрыть почти все ненужные приложения, перформанс моего устройства оставляет желать лучшего ввиду его слабеньких комплектующих

### Результаты 

Ниже я прикреплю несколько графиков, описывающих результаты тестирования

![First results](/media/media_1.png)

Данный график показывает резултаты тестирования во время линейного роста нагрузки до ~ 850 RPS. Медианное время ответа при такой нагрузке составило 20-30 мс, при этом все запросы успешно получили ответ

На следующем участке система начала голодать и терять стабильность, но количество ошибок не превышало 10. Из-за голодания, задержек и, вероятно, оверхеда от смена контекстов процессора медианное время ответо начало тянуться сначала к 80-100 мс, а потом к 400 мс. Это вы можете увидеть в правой половине следующих графиков:

![Second results](/media/media_2.png)

Тестирующая система несколько раз выдавала варнинги о высоком потреблении процессора и нехватке мощностей, поэтому я попробовал запустить её в 4 воркера (и забыл закрыть VSCode), что привело к голоданию приложения и оно не смогло заскейлиться выше 400 RPS, а под конец начало возвращать нулевой статус ответа на один из эндпоинтов:

![Failed results](/media/media_3.png)

### Интерпретация результатов

Не считая скачков времени ответа, система показала низкое количество ошибок и неплохо справилась с тестированием. Одна из возможных причин таких скачков (Но это неточно) - CPU голодание самой тестирующей системы, которая могла иногда не справляться с генерацией такой нагрузки. Я предполагаю, что камнем преткновения в приложении могла стать база данных. В будущем для выдерживания большей нагрузки можно распилить систему на микросервисы, или хотя бы сделать бд распределенной. Я также не знаю насколько хорошо с подобной нагрузкой справляется chi-router и стандартная библиотекa, поэтому возможно стоило бы рассмотреть другой фреймворк. Если же проблема не в базе данных, то можно просто размножить инстансы приложения на разных устройствах и балансировать нагрузку по ним, например с помощью доменных имён или использования reverse proxy.