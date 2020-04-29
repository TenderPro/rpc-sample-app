# rpc-sample-app

## Описание
Пример (для fork) cервиса на golang, демонстрирующий использование технологий:

* gRPC
* OpenAPI
* Websocket
* SOAP
* NRPC (gRPC via NATS)
* NATS pub/sub

в котором весь обмен данными между подсистемами, включая все структуры запросов, ответов и коды ошибок, описан в [protobuf](main.proto)

## Варианты работы сервиса

1. **mono** - полный функционал без NRPC (monolyth: proxy <-> handler)
2. **bus** - полный функционал с обменом по NRPC между gRPC-proxy и сервером (proxy <-> bus <-> handler)
3. **handler** - функционал сервера, который подключается к NATS в качестве сервера
4. **proxy** - функционал gRPC-proxy, который подключается к NATS в качестве клиента

## Структура приложения

```
├── docker-compose.yml
├── Dockerfile
├── go.mod - в текущей версии требует установки в соседнем каталоге [rpckit](github.com/TenderPro/rpckit)
├── main.go
├── main.proto - описание сервиса
├── Makefile
├── pkg
│   ├── app - сборка сервиса из компонентов
│   ├── nrpcgen
│   ├── pb
│   ├── service - прикладная часть сервиса
│   ├── soapgen
│   ├── staticgen
│   └── template - работа с шаблонами
├── README.md
└── static
    ├── html
    ├── sql
    └── tmpl

```

## Описание реализации

Для сборки проекта используется go версии 1.14. По описанию сервиса, [main.proto](main.proto), автогенерится код поддержки NATS-RPC, OpenAPI, SOAP документация (в каталоги pkg/*gen и pkg/pb). Содержимое каталога `static` командой `make gen-prod` переностится в пакет `pkg/staticgen`.

Команда сборки сервиса: `go build .`

Сервис поддерживает следующие аргументы командной строки:
```
$ ./rpc-sample-app -h

Usage:

```

### Использование make

Для облегчения повторного запуска можно использовать команды `make`:
* `make conf` - создать файл конфигурации .env
* `make run` - локальная сборка и запуск сервиса с конфигурацией из .env

Полный список команд:
```
$ make help
...
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

Файл [server_test.go](server_test.go) позволяет провести тестирование методов API при работающем сервисе, но текущая версия не удаляет после запуска изменения в БД

## Дополнения

* golint выдает замечания на использование `Id`, но такая особенность protobuf [документирована](https://github.com/golang/protobuf/issues/73#issuecomment-138699104)
* protoc добавляет в структуры поля с префиксом `XXX_`, что мешает использовать эти структуры в gorm. Для решения можно было бы [использовать gogo/protobuf](https://github.com/golang/protobuf/issues/52#issuecomment-284219742) или [retag](https://github.com/golang/protobuf/issues/52#issuecomment-295596893), но пришлось бы добавить комменты в api.proto

## TODO

* [x] ping.Timeservice
* [x] актуализировать примеры в static/html и тесты в Makefile
* [ ] актуализировать README
* [ ] pkg/app.Run - что еще вынести в rpckit?
* [ ] пример вызова метода из шаблонов
* [ ] gRPC: возврат ошибок
* [ ] nrpc: трейсинг
* [ ] nrpc: protoc-gen (fork)
* [ ] pgmig: пример работы с БД
* [ ] pgmig: protoc-gen
* [ ] пример file upload (sfs)
* [ ] make lint
* [ ] make cov (часть 1 - корректная работа)
* [ ] make cov (часть 2 - тесты с docker)
* [ ] make cov (часть 3 - coverage >80%)
* [ ] https://codebeat.co/
* [ ] доработать документацию

## License

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2020 Tender.Pro <it@tender.pro>
Copyright (c) 2019 Aleksei Kovrizhkin <lekovr@gmail.com>
