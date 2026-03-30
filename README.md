# P2P gRPC Chat (Go)

Консольный peer-to-peer чат на Go.

## Тестирование

- Запуск всех тестов: `go test ./...`
- Запуск тестов через `make`: `make test`
- Просмотр покрытия: `make cover`

Команда `make cover`:

- запускает все тесты с `coverprofile`;
- генерирует файл `coverage.html`;
- сразу открывает HTML-отчёт в браузере.

Подробный план тестирования и чеклист реализованных тестов находятся в [docs/testing.md](docs/testing.md).
