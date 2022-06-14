# task rest api

## Общий функционал

rest api
```
POST   /tasks/              :  создаёт задачу и возвращает её 
GET    /tasks/              :  возвращает все задачи

DELETE /tasks/<taskid>      :  удаляет задачу по ID
DELETE /tasks/              :  удаляет все задачи

GET    /tasks/<taskid>      :  возвращает одну задачу по её ID
GET    /tasks/?tag=<tag>    :  возвращает список задач с заданным тегом
GET    /tasks/?date=<date> :  возвращает список задач, запланированных на указанную дату  формат даты "2006-01-02"
```

## Запуск

```bash
export LogLevel=10
export Port=8082

# создание sqlite db
sqlite3 test.db < scripts/create_db.sql

# сборка 
go build  -o . ./...
```

## Основные технологии

- http 
- pat роутер
- aconfig - парсер конфигов 
- alice - chunk middleware
- go-sqlite3
- sync.errgroup


## Структура приложения

- internal/rest - http сервер, middleware, роуты
- internal/logger - логгер
- internal/models - модели, интерфейс репозитория
- internal/taskstore/sqlstore - общая реализация стора для sql
- internal/taskstore/sqlitestore - реализация стора для sqlite3

## TODO 

1. упаковать в docker compose
2. добавить тесты
3. добавить враппер ErrorWrapper для всех ошибок по аналогии с TaskNotFound
3. баг - GetAllTags не возвращает теги
