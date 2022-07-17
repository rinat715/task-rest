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

POST /users                 : создает юзера
GET /users/<userid>         : возвращает юзера по ID
```

## Запуск

```bash
make build
make init
```

## Основные технологии

- http 
- pat роутер
- aconfig - парсер конфигов 
- alice - chunk middleware
- go-sqlite3
- go migrate
- easyjson
- validator
- wire 


## Структура приложения

- internal/rest - http сервер, middleware, роуты
- internal/logger - логгер
- internal/models - модели, интерфейс репозитория
- internal/repositories - репозитории
- internal/services

## TODO 

1. добавить враппер ErrorWrapper для всех ошибок по аналогии с TaskNotFound
