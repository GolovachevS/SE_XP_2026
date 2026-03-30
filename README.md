# P2P gRPC Chat (Go)

[![CI](https://github.com/GolovachevS/SE_XP_2026/actions/workflows/ci.yml/badge.svg)](https://github.com/GolovachevS/SE_XP_2026/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Консольный peer-to-peer чат на Go с прямым соединением между двумя узлами через gRPC bidirectional streaming.

## Команда

- Артамонова Анна 
- Головачев Сергей

## Возможности

- прямое peer-to-peer соединение без центрального сервера;
- консольный интерфейс;
- отображение имени отправителя, времени и текста сообщения;
- работа в двух режимах:
  - ожидание входящего подключения;
  - подключение к другому peer;
- graceful shutdown при завершении сессии и отмене приложения;
- автоматические unit- и integration-тесты.

## Требования

- Go `1.26`
- `protoc`
- плагины:
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`

Примечание: protobuf/gRPC stubs уже лежат в репозитории. `make proto` нужен, если меняется файл `api/chat/v1/chat.proto`.

## Параметры запуска

Приложение запускается через `cmd/chat`.

Поддерживаемые флаги:

- `-name <username>`: имя пользователя, обязательный параметр;
- `-peer <host:port>`: адрес peer-а для исходящего подключения;
- `-listen <host:port>`: адрес и порт для ожидания входящего подключения, по умолчанию `:50051`.

Если `-peer` не указан, приложение работает в режиме ожидания входящего подключения.

## Быстрый запуск

Сборка:

```bash
make build
```

Запуск сервера:

```bash
./bin/chat -name Alice -listen :50051
```

Запуск клиента:

```bash
./bin/chat -name Bob -peer 127.0.0.1:50051
```

Пример сценария в двух терминалах:

```bash
# Терминал 1
./bin/chat -name Alice -listen :50051

# Терминал 2
./bin/chat -name Bob -peer 127.0.0.1:50051
```

После установления соединения обе стороны могут отправлять сообщения друг другу.

## Основные команды

- `make build` — собрать бинарник `bin/chat`
- `make test` — запустить все тесты
- `make cover` — построить coverage report и открыть `coverage.html`
- `make fmt` — отформатировать код
- `make fmt-check` — проверить форматирование
- `make lint` — запустить `golangci-lint`
- `make proto` — пересгенерировать protobuf/gRPC код

Также можно запускать тесты напрямую:

```bash
go test ./...
```

## Структура проекта

- `cmd/chat` — точка входа и запуск приложения;
- `internal/app` — orchestration-слой, режимы работы и session loop;
- `internal/chat` — доменная модель сообщения;
- `internal/transport/grpcchat` — реализация транспорта поверх gRPC;
- `internal/ui/console` — консольный ввод/вывод;
- `internal/protocol/chatv1` — mapping между domain model и protobuf;
- `api/chat/v1` — protobuf-контракт и сгенерированный код;
- `docs/architecture.md` — архитектурная документация;
- `docs/testing.md` — план тестирования и чеклист;
- `docs/task-decomposition.md` — декомпозиция задач в команде.

## Тестирование и отчёты

- план тестирования: [docs/testing.md](docs/testing.md)
- архитектурная документация: [docs/architecture.md](docs/architecture.md)
- декомпозиция задач: [docs/task-decomposition.md](docs/task-decomposition.md)
- покрытие: `make cover` создаёт `coverage.out` и `coverage.html`

## Важные замечания

- поддерживается сценарий `1↔1`;
- адрес peer-а задаётся вручную;
- не реализуются NAT traversal и relay;
- `client/server` в проекте означают только роли peer-ов при установлении соединения, а не наличие центрального сервера.

## Лицензия

Проект распространяется под лицензией MIT. См. [LICENSE](LICENSE).
