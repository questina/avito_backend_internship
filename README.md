# avito_backend_internship
Микросервис для работы с балансом пользователей

Структура проекта:
```
.
├── Dockerfile
├── README.md
├── db
│   └── db.go
├── docker-compose.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── init.sql
├── reports
│   └── report_11_2022.csv
├── schemes
│   └── structs.go
├── server.go
├── server_test.go
└── utils
    ├── routes-handlers.go
    ├── routes.go
    ├── start_init.go
    └── utils.go
```
Архитектурные решения:
* В основной папке лежит Dockerfile и docker-compose для создания контейнера с сервисом
* ./server.go - пакет main
* ./server_test.go - тестировщик сервиса, успела написать тесты только на основные ручки сервиса
* ./init.sql - файл для создания базы данных и таблиц
* ./db - содержит в себе функцию для подключения к БД MySQL
* ./schemes - содержит основные структуры проекта для работы с json, а также для удобства работы с сервисом
* ./utils - основные функции. Добавляются и реализуются ручки, запросы к БД, старт сервера
* ./docs - файлы, сгенерированные swag для языка Golang. Поднять API не получилось, но есть документ сваггера.
* ./reports - файл с бухгалтерскими отчетами для дополнительного задания

Запустить сервис:
```bash
go run .
```
Или запустить через докер:
```bash
docker-compose up
```

Запустить тесты:
```bash
go test
```

## Все основные запросы через CURL (можно добавить флаг -v)

Для создания пользователя нужно зачислить деньги на счет. Система вернет id созданного пользователя и запишет его в БД.
```bash
curl -X POST http://localhost:8000/add_money -H 'Content-Type: application/json' -d '{"amount":100}'
```
Должно получится что-то такое:
```bash
{"Id":1,"Balance":100,"Status":"OK"}
```
Чтобы добавить деньги на счет уже существующего пользователя достаточно передать его id через json.
```bash
curl -X POST http://localhost:8000/add_money -H 'Content-Type: application/json' -d '{"amount":99.9,"id":1}'
```
Баланс должен увеличиться на 99.9
```bash
{"Id":1,"Balance":199.9,"Status":"OK"}
```
Далее можно резервировать деньги на балансе пользователя. Для этого нужно передать id пользователя, id его заказа, id сервиса, который обрабатывает заказ и стоимость заказа. 
```bash
curl -X POST http://localhost:8000/reserve_money -H 'Content-Type: application/json' -d '{"userid":1,"orderid":1,"serviceid":1,"cost":10}'
```
Должен вернуться статус ОК.
Проверяем что деньги зарезервировались
```bash
curl -X POST http://localhost:8000/get_balance -H 'Content-Type: application/json' -d '{"id":1}'
```
Получаем: 
```bash
{"Id":1,"Balance":189.9}
```
Если все ок, то можно забрать зарезервированные деньги. Передаем такие же параметры.
```bash
curl -X POST http://localhost:8000/take_money -H 'Content-Type: application/json' -d '{"userid":1,"orderid":1,"serviceid":1,"cost":10}'
```
Должен вернуться статус ОК.

Если возникла какая-то проблема, то деньги можно аналогично разрезервировать.
```bash
curl -X POST http://localhost:8000/free_money -H 'Content-Type: application/json' -d '{"userid":1,"orderid":1,"serviceid":1,"cost":10}'
```
Проверим баланс пользователя:
```bash
curl -X POST http://localhost:8000/get_balance -H 'Content-Type: application/json' -d '{"id":1}'
```
Получаем: 
```bash
{"Id":1,"Balance":199.9}
```

Также программа может создавать бухгалетрские отчеты, указывая год и месяц за который они нужны.
```bash
curl -X POST http://localhost:8000/generate_report -H 'Content-Type: application/json' -d '{"Month":11,"Year":2022} --output tmp.csv'
```
Результат запишется в файл tmp.csv (можно любой другой файл). 

Пользователи могут также посмотреть информацию о своих деньгах (где и когда они были получены или потрачены). Можно также указать query параметры (сортировка и limit, offset для пагинации)
```bash
curl -X POST http://localhost:8000/balance_info?sort=desc&limit=5&offset=0 -H 'Content-Type: application/json' -d '{"id":1}'
```
Получим:
```bash
[{"Timestamp":"2022-11-09 20:27:58","Amount":10,"EventType":"FREE","ServiceId":1,"OrderId":1},
{"Timestamp":"2022-11-09 20:21:07","Amount":10,"EventType":"RESERVE","ServiceId":1,"OrderId":1},
{"Timestamp":"2022-11-09 20:16:57","Amount":10,"EventType":"RESERVE","ServiceId":1,"OrderId":1},
{"Timestamp":"2022-11-09 20:14:30","Amount":100,"EventType":"ADD","ServiceId":1,"OrderId":1}]
```
### Далее напишу тактические решения:
База данных содержит три таблицы (подробнее см. init.sql): 
- user_balances(user_id, balance, reserved): содержит информацию о пользователе, его балансе и сколько денег у него зарезервировано
- orders (order_id, service_id, user_id, cost): содержит информацию о заказах. Считаю что order_id уникален.
- moneyflow (datetime, event_type, order_id, service_id, user_id, amount): содержит информацию о всех действиях с сервисом

Ручка **/add_money** добавляет новую запись в таблицу *user_balances* или обновляет баланс старой. Если вводят отрицательную сумму денег, возникает ошибка. 

Ручка **/reserve_money** проверяет что такой *order_id* не содержится в таблице *orders*, а также что на счету пользователя достаточно средств для этой операции с учетом зарезервированных денег. Далее я делаю транзакцию, которая состоит из двух операций с БД: добавление новой записи в таблицу *orders*, а также обновление поля *reserved* в строке пользователя в таблице *user_balances*. Если возникает какая-то ошибка в одной из этих операций, то транзакция не выполняется и база данных возвращается к изначальному состояния. 

Ручка **/take_money** проверяет что order_id есть в таблице *orders*. Если все ок, то аналогично выполняется транзакция из двух операций: удаление записи о резервированных деньгах и обновление записи в таблице пользователей об его зарезервированных деньгах. 

Ручка **/get_balance** возвращает текущий баланс пользователя (с учетом зарезервированных денег)

Ручка **/generate_report** выполняет запрос к таблице *moneyflow* и создает локальный файл с отчетом в папке *reports* и отправляет его клиенту.

Ручка **/balance_info** выполняет запрос к таблице *moneyflow* со всеми заданными параметрами и отправляет json с событиями пользователю. 

Тесты не успеваю дописать к сожалению, там только на первые четыре ручки и с основными сценариями
