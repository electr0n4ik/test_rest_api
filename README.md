# Запуск Проекта:

### Установка зависимостей:

```shell
go get -u
```

### Инициализация модуля Go:

```shell
go mod init test_rest_api
```

### Запуск сервера:

```shell
go run main.go.
```

### JSON для Postman:

    {
    "title": "Album 1",
    "artist": "Artist 1",
    "price": 29.99
    }

### Создание записи без Postman:

```shell
curl http://localhost:8080/albums \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"id": 4, "title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
```

### Получить список созданных записей:

```shell
curl http://localhost:8080/albums \
--header "Content-Type: application/json" \
--request "GET"
```
