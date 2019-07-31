# companyserv

## 1. Описание

Сервис на golang, предоставляющий gRPC API. API сервиса, включая все структуры запросов, ответов и коды ошибок описаны в [protobuf](api.proto)

## 2. Требования к Сервису

1. сервис должен работать в отдельной БД (Postgresql) с заданной структурой таблиц (схема данных - [companyserv.sql](companyserv.sql)
2. сервис должен конфигурироваться (хост бд, пользователь, адрес интерфейса gRPC и т.д.) через задание опций командной строки. 
3. методы GetCUserPolicies и UpdateCUserPolicies реализовать заглушками, возвращающими всегда пустой ответ без ошибки.
4. SwitchCUserStatus определяет допустимые переходы по таблице status_transitions



## 3. Нюансы

* в companyserv.sql зашито имя пользователя (ma), оно используется по умолчанию. Если его изменить в настройках, разворачивание БД повлечет ошибки
* формулировка `X СОДЕРЖИТ указанный текст`: буквально это `LIKE %X%`, но в большинстве случаев лучше `~* X`, в решении использован первый вариант
* GetCompanyIDs, GetCUserIDs: нумерация страниц начинается с 0
* golint выдает замечания на использование `Id`, но такая особенность protobuf [документирована](https://github.com/golang/protobuf/issues/73#issuecomment-138699104)
* protoc создает в структурах поля с префиксом `XXX_`, что мешает использовать эти стрктуры в gorm. Для решения можно было бы [использовать gogo/protobuf](https://github.com/golang/protobuf/issues/52#issuecomment-284219742) или [retag](https://github.com/golang/protobuf/issues/52#issuecomment-295596893), но пришлось бы добаивть комменты в api.proto
* Ошибка внешнего ключа для PG возвращается в таком виде: `Severity:"ERROR", Code:"23503", Message:"insert or update on table \"company_users\" violates foreign key constraint \"fk_rails_946619ff40\"", Detail:"Key (user_id)=(-2) is not present in table \"users\".", Table:"company_users", Constraint:"fk_rails_946619ff40"`, если ключи называть осмысленно, по имени можно было бы сказать, с какой таблицей связана ошибка. Без этого приходится проверять отдельно и в транзакции

## 4. Расширения

* в .proto нет метода создания записи в user, решено в addon.sql


* [ ] тесты: режим одной транзакции
* [ ] тесты grpc
* [ ] моки?
