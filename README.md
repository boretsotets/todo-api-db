# Task tracker API

Простой API для управления задачами.
Стек: Go, PostgreSQL, JWT авторизация, миграции через [golang-migrate](https://github.com/golang-migrate/migrate).
Запуск и работа в docker контейнерах. Реализован по принципам чистой архитектуры

## Быстрый старт

```bash
docker-compose up --build
```
API будет доступно на http://localhost:8080

## Возможности

- Регистрация и авторизация пользователей
- CRUD для задач
- Поддержка миграций БД

## Примеры запросов

```
# пример регистрации. Нужно указать поля "name", "email" и "password".
curl -X POST -H "Content-Type application/json" -d '{"name": "username", "email": "user@name.com", "password": "password1"}' http://localhost:8080/register

# пример авторизации. Нужно указать поля "email" и "password".
curl -X POST -H "Content-Type application/json" -d '{"email": "user@name.com", "password": "password1"}' http://localhost:8080/login

# пример создания задачи. Нужно указать поля "title" и "description".
curl -X POST -H "Content-Type application/json" -H "Authorization: ..." -d '{"title": "title1", "description": "description1"}' http://localhost:8080/todos

# пример изменения задачи. Нужно указать поля "title" и "description" и указать id изменяемой задачи в URL.
curl -X PUT -H "Content-Type application/json" -H "Authorization: ..." -d '{"title": "newtitle", "description": "newdescription"}' http://localhost:8080/todos/1

# пример получения списка задач. Нужно указать номер страницы и размер страницы в URL.
curl -X GET -H "Content-Type application/json" -H "Authorization: ..." "http://localhost:8080/todos?page=1&limit=10"

# пример удаления задачи. Нужно указать id удаляемой страницы в URL.
curl -X DELETE -H "Content-Type application/json" -H "Authorization: ..." http://localhost:8080/todos/1
```
