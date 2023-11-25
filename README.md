# Запуск Проекта:

Выполнение миграций:
```shell
migrate -database postgres://postgres:12345@localhost:5432/postgres?sslmode=disable -path db/migrations up
```

```shell
go run main.go.
```