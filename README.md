![CI](https://github.com/Ekstrem/DigiTFactory.Go.Libraries/actions/workflows/ci.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/Ekstrem/DigiTFactory.Go.Libraries.svg)](https://pkg.go.dev/github.com/Ekstrem/DigiTFactory.Go.Libraries)

# DigiTFactory.Go.Libraries

Фреймворк быстрой разработки микросервисов, декомпозированных по субдомену.
Go-порт [DigiTFactory.Libraries](https://github.com/Ekstrem/DigiTFactory.Libraries) (.NET).

## Проблема

При производстве программного обеспечения на основе архитектурного стиля микросервисы очень важным выступает предметно-ориентированное проектирование. У разработчиков очень часто возникают проблемы обучения, фокусировки на тактических паттернах и в принципе чисто технических проблемах. При этом теряется самое главное — моделирование, коммуникация, фиксация на главных задачах бизнеса.

Цель этого проекта — дать продуктовым командам готовую инфраструктуру для разработки микросервисов, с помощью которой разработчики смогут начать разработку быстрее, одновременно проводя своё обучение.

**Никаких тренингов, инфобизнеса, деформаций и дисфункций** — проектируйте и исследуйте!

## Установка

```bash
go get github.com/Ekstrem/DigiTFactory.Go.Libraries@latest
```

## Архитектура

Гексагональная архитектура. Микросервис представлен следующими слоями:

1. **Доменная модель** — агрегаты, объекты-значения, спецификации
2. **Доменные сервисы** — провайдер агрегатов, нотификатор шины событий
3. **Уровень приложения** — обработчики команд и запросов
4. **Уровень инфраструктуры** — реализации репозиториев, шин событий
5. **Уровень хостинга** — HTTP/gRPC серверы

### Строение модели

| Термин | Описание |
|---|---|
| **Подобласть** | Идеально выделенный экспертами предметной области ограниченный контекст |
| **Ограниченный контекст** | Содержит агрегаты и доменные сообщения |
| **Агрегат** | Анемичная модель + границы области (валидаторы и бизнес-операции) |
| **Корень агрегата** | Версионированная сущность — одновременно корневой value object и entity |
| **Бизнес-операция** | Объект-фабрика: валидирует модель → создаёт новую версию агрегата → проецирует доменное событие |

### Важные особенности

1. **Иммутабельность**. Все типы модели намеренно иммутабельны. Изменение агрегатов проводится только через бизнес-операцию, что гарантирует инвариант агрегата и границы транзакционной целостности.
2. **Event Bus в центре**. База данных — лишь адаптер. Центром всех команд является шина доменных событий.

## Состав пакетов

### `seedworks/` — Ядро DDD

Основные тактические паттерны и контракты.

```
seedworks/
├── tacticalpatterns/   # Aggregate, Entity, ValueObject, BusinessOperation, ReadModel
├── characteristics/    # HasKey, HasVersion, HasParent, Paging, CorrelationToken
├── definition/         # BoundedContext, Interaction
├── events/             # DomainEvent, EventBus, EventHandler, Notifier
├── invariants/         # Specification, Validator, ValidateCommand
├── repository/         # CommandRepository, QueryRepository, ReadModelStore, UnitOfWork
├── result/             # AggregateResult, OperationResult, TaskResult
├── monads/             # Result[T], Extensions
├── lifecycle/          # ValueObjectHelper
└── reactive/           # Unsubscriber
```

### `eventbus/` — Реализации шины событий

| Пакет | Описание | Зависимость |
|---|---|---|
| `eventbus/inmemory` | In-process, для тестов и монолитов | — |
| `eventbus/kafka` | Продакшн, распределённая шина | `segmentio/kafka-go` |
| `eventbus/postgres` | Outbox-паттерн, транзакционная надёжность | `jackc/pgx/v5` |

### `commandrepository/` — Event Store (запись)

Стратегии хранения событий (`commandrepository/strategy.go`):

| Стратегия | Описание |
|---|---|
| `FullEventSourcing` | Все события хранятся, агрегат восстанавливается из полного потока |
| `SnapshotAfterN` | Снимок каждые N событий, чтение: снимок + последующие события |
| `StateOnly` | Только текущее состояние агрегата (без истории событий) |

| Пакет | Хранилище | Зависимость |
|---|---|---|
| `commandrepository/postgres` | PostgreSQL | `jackc/pgx/v5` |
| `commandrepository/mongo` | MongoDB | `mongo-driver` |

### `queryrepository/` — Read Store (чтение)

| Пакет | Хранилище | Зависимость |
|---|---|---|
| `queryrepository/postgres` | PostgreSQL | `jackc/pgx/v5` |
| `queryrepository/redis` | Redis | `go-redis/v9` |
| `queryrepository/scylla` | ScyllaDB / Cassandra | `gocql` |

## Быстрый старт

### Определение доменного события и обработчика

```go
package main

import (
    "context"
    "log/slog"

    "github.com/Ekstrem/DigiTFactory.Go.Libraries/seedworks/events"
    "github.com/Ekstrem/DigiTFactory.Go.Libraries/eventbus/inmemory"
)

type OrderHandler struct{}

func (h *OrderHandler) HandleEvent(ctx context.Context, event events.DomainEvent) error {
    slog.Info("order event received",
        "aggregate", event.AggregateID,
        "version", event.Version,
    )
    return nil
}

func main() {
    bus := inmemory.New(slog.Default())

    handler := &OrderHandler{}
    bus.Subscribe("OrderCreated", handler)

    event := events.NewDomainEvent(/* ... */)
    bus.Publish(context.Background(), event)
}
```

## CI/CD

Pipeline запускается автоматически на push в `main` и при создании PR:

1. **Build & Test** — `go build ./...` + `go test -race ./...`
2. **Lint** — `golangci-lint`
3. **Release** — при пуше тега `v*` модуль индексируется на `proxy.golang.org`

### Публикация новой версии

```bash
git tag v0.2.0
git push --tags
```

Модуль автоматически появится на [pkg.go.dev](https://pkg.go.dev/github.com/Ekstrem/DigiTFactory.Go.Libraries).

## Зависимости

| Пакет | Версия | Назначение |
|---|---|---|
| `google/uuid` | v1.6.0 | Генерация UUID |
| `jackc/pgx/v5` | v5.7.2 | PostgreSQL драйвер |
| `segmentio/kafka-go` | v0.4.47 | Kafka клиент |
| `mongo-driver` | v1.17.2 | MongoDB драйвер |
| `go-redis/v9` | v9.7.0 | Redis клиент |
| `gocql` | v1.7.0 | ScyllaDB/Cassandra драйвер |

## Лицензия

См. [LICENSE](LICENSE).
