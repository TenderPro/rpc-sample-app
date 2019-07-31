# companyserv

## Описание

Сервис на golang, предоставляющий gRPC API. API сервиса, включая все структуры запросов, ответов и коды ошибок описаны в [protobuf](grpccompanyserv/api.proto)

## Требования к Сервису

1. сервис должен работать в отдельной БД (Postgresql) с заданной структурой таблиц (схема данных - [companyserv.sql](sql/companyserv.sql)
2. сервис должен конфигурироваться (хост бд, пользователь, адрес интерфейса gRPC и т.д.) через задание опций командной строки. 
3. методы GetCUserPolicies и UpdateCUserPolicies реализовать заглушками, возвращающими всегда пустой ответ без ошибки.
4. SwitchCUserStatus определяет допустимые переходы по таблице status_transitions

## Описание реализации

Для сборки проекта использовался go версии 1.12.6. Результат вызова `protoc`, файл api.pb.go, включен в проект и эта зависимость будет нужна только в случае изменения api.proto.

Команда сборки сервиса: `go build .`

Сервис поддерживает следующие аргументы командной строки:
```
$ ./ grpcsample-h

Usage:
   grpcsample[OPTIONS]

Application Options:
      --addr=         Listen address (default: localhost:7070)
      --debug         Print debug logs

API Options:
      --api.maxcount= Id slice max len (default: 1000)

DB Options:
      --db.addr=      host:port (default: localhost:5432)
      --db.driver=    DB driver (default: postgres)
      --db.user=      User name
      --db.password=  User password
      --db.name=      Database name
      --db.opts=      Database connect options (default: sslmode=disable)

Help Options:
  -h, --help          Show this help message

```

### Использование make

Для облегчения повторного запуска можно использовать команды `make`:
* `make conf` - создать файл конфигурации .env
* `make run` - локальная сборка и запуск сервиса с конфигурацией из .env

Полный список команд:
```
$ make help
api                            Generate grpc go sources
build                          Build the binary file for server
clean                          Remove previous builds
conf                           Create initial config
dcape-db-create                Create user, db and load dump
dcape-db-drop                  Drop database and user
dcape-psql                     Run psql
dc                             Run docker-compose (make dc CMD=build)
dep                            Get the dependencies
down                           Stop containers and remove them
help                           Display this help screen
lint                           Run linter
psql                           Run psql via postgresql docker container
run                            Build and run binary
test                           Run grpc client tests
up-db                          Start pg container only
up                             Start pg and app containers

```

### Использование docker

При наличии локально установленных make и docker, сборка и запуск сервиса могут быть произведены командой
```
make up
```
Выполнение этой команды повлечет
* запуск локальной копии БД
* загрузку в БД файлов из sql
* сборку проекта
* запуск сервиса проекта

при этом будут использованы образы docker:
* docker/compose:1.23.2
* golang:1.12.6-alpine3.9
* postgres:11.4

### Использование dcape

Сервис также поддерживает деплой в рамках сервиса [dcape](https://github.com/dopos/dcape)

### Тесты

В дополнение к запрошенному в задаче функционалу, файл [server_test.go](server_test.go) позволяет провести тестирование методов API при работающем сервисе, но текущая версия не удаляет после запуска изменения в БД

## Дополнения

* в companyserv.sql зашито имя пользователя (ma), оно используется по умолчанию. Если его изменить в настройках, разворачивание БД повлечет ошибки
* формулировка `X СОДЕРЖИТ указанный текст`: буквально это `LIKE %X%`, но в обычно используют `~* X`, в решении выбран первый вариант
* GetCompanyIDs, GetCUserIDs: нумерация страниц начинается с 0
* golint выдает замечания на использование `Id`, но такая особенность protobuf [документирована](https://github.com/golang/protobuf/issues/73#issuecomment-138699104)
* protoc добавляет в структуры поля с префиксом `XXX_`, что мешает использовать эти структуры в gorm. Для решения можно было бы [использовать gogo/protobuf](https://github.com/golang/protobuf/issues/52#issuecomment-284219742) или [retag](https://github.com/golang/protobuf/issues/52#issuecomment-295596893), но пришлось бы добавить комменты в api.proto
* Ошибка внешнего ключа для PG возвращается в таком виде: `Severity:"ERROR", Code:"23503", Message:"insert or update on table \"company_users\" violates foreign key constraint \"fk_rails_946619ff40\"", Detail:"Key (user_id)=(-2) is not present in table \"users\".", Table:"company_users", Constraint:"fk_rails_946619ff40"`, если ключи называть осмысленно, по имени можно было бы сказать, с какой таблицей связана ошибка. Без этого приходится проверять отдельно и в транзакции
* файл [sql/z_addon.sql](sql/z_addon.sql) содержит дополнения в БД, необходимые для запуска тестов

## Автор

2019, Алексей Коврижкин <lekovr@gmail.com>
