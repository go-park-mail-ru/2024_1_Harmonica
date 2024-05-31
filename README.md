# Harmonica 2024
### Участники команды
* Купцов Гавриил
* Костин Глеб
* Амирова Лилиана
### [Frontend repository](https://github.com/frontend-park-mail-ru/2024_1_Harmonica)
### Deploy: [harmoniums.ru](https://harmoniums.ru/)
### [Figma](https://www.figma.com/design/zRx9iBFVMZe01acfiQyfzO/My-Pinterest?node-id=0-1&t=4PLYeUGCXQmhgaCB-0)
### [Swagger](https://harmoniums.ru:8080/swagger/)
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

В любом доступном интерфейсе (например, [PgAdmin](https://www.pgadmin.org/download/) или [DBevear](https://dbeaver.io/download/)) подключитесь к созданной базе и выполните SQL скрипты, описанные в [файле](../main/db/migrations/initDB.sql).
#### 2. Установка S3 Minio
Установите [Minio](https://min.io/). 

Необходимо поднять сервер minio, создать bucket "images" и поместить следующие параметры в conf.env (создайте файл в корне проекта):
```
MinioEndpoint=
MinioAccessKeyID=
MinioSecretAccessKey=
MinioUseSSL=false
```

#### 3. Подъем сервера.
В корень проекта поместите файл **conf.env**, добавив туда содержимое: 
```env
# DB Config
DBHost=localhost
DBPort=5432
DBUser=postgres
DBPassword=postgres
DBname=harmonica
SERVER_URL=http://127.0.0.1:8080/img/
AUTH_MICROSERVICE_PORT=:8002
IMAGE_MICROSERVICE_PORT=:8003
LIKE_MICROSERVICE_PORT=:8004
DEBUG=true
```
Чтобы запустить сервер в консоли, перейдите в корень проекта и запустите следующую команду:

`make run`

