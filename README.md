# 2024_1_Harmonica
### Участники команды
* Купцов Гавриил
* Костин Глеб
* Амирова Лилиана
### [Frontend repository](https://github.com/frontend-park-mail-ru/2024_1_Harmonica)
---
### Инструкция по локальной сборке
#### 1. Установка БД PostgreSQL
Для установки Postgres используется одна из следующих команд, в зависимости от ОС: 
* MacOS:
`brew install postgresql`
* Linux:
`sudo apt-get install postgresql postgresql-contrib`
* Windows:
`runas /user:postgres cmd.exe`

Необходимо войти в интерфейс postgres:
`psql postgres (для выхода из интерфейса используйте \q)`

Затем создайте и настройте локальную БД, использовав команды:
```
CREATE DATABASE pinterest;
CREATE USER postgres;
ALTER ROLE postgres SUPERUSER PASSWORD "postgres";
```

В любом доступном интерфейсе (например, [PgAdmin](https://www.pgadmin.org/download/) или [DBevear](https://dbeaver.io/download/)) подключитесь к созданной базе и выполните SQL скрипты, описанные в [файле](../db/migrations/initDB.sql).
#### 2. Подъем сервера.
В корень проекта поместите файл **conf.go**, скопировав туда содержимое: 
```go
package main

import "harmonica/db"

var Conf = db.DBConf{
	Host:     "localhost",
	Port:     5432,
	User:     "postgres",
	Password: "postgres",
	Dbname:   "pinterest",
}
```
Чтобы запустить сервер в консоли, перейдите в корень проекта и запустите следующую команду:

`go run harmonica`
