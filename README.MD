Решение тестового задания на позицию стажера-бекендера
https://github.com/avito-tech/autumn-2021-intern-assignment
Микросервис для работы с балансом пользователей

Для обеспечения работы микросервисы была выбрана база данных PostgreSQL
Для развертываения БД принято решение использовать docker
Команда: docker run --name avitoms -p 5432:5432 -e POSTGRES_USER=avitoms -e POSTGRES_PASSWORD=password -e POSTGRES_DB=avitoms -d postgres:13.3

Для работы с базой данных использован пакет github.com/jinzhu/gorm

Данные для подключения к БД хранятся в файле .env
Для чтения .env использован пакет github.com/joho/godotenv

API
/api/account/new GET создает новый аккаунт
возвращает id созданного аккауна

/api/account/balance GET проверка баланса
принимает 
user_id - id аккаунта
currency - Валюта в которой нужно вывести баланс. Если не указано баланс по умолчанию выводится в рублях
возвращает 
balance - баланс аккаунта

/api/account/credit GET зачисление средств
принимает
user_id - id аккаунта
sum - сумма
возвращает
transact_id - id транзакции

/api/account/debit GET списание средств
принимает
user_id - id аккаунта
sum - сумма
targer - цель списания средств (корзина, товар...)
возвращает
transact_id - id транзакции

/api/account/transfer GET перевод средств между аккаунтами
принимает
user_id -id аккаунта с которого списываются деньги
user_id_to - id аккауна на который зачисляются деньги
sum - сумма
возвращает
transact_id_user - id транзакции для аккаунта с которого списывают деньги
transact_id_user_to - id транзакции для аккаунта которому зачисляются деньги

/api/account/transacts GET вывод транзакций
принимает
user_id - id аккаунта
возвращает список транзакций

Примеры postman
создаем два аккаунта
localhost:8000/api/account/new (id:1)
localhost:8000/api/account/new (id:2)
зачисляем деньги на баланс аккауна id=1
localhost:8000/api/account/credit?user_id=1&sum=3000
списываем деньги на товар для аккаунта id=1
localhost:8000/api/account/debit?user_id=1&sum=600&target=iphone
переводим деньги от аккауна id=1 к id=2
localhost:8000/api/account/transfer?user_id=1&user_id_to=2&sum=300
проверяем баланс аккауна id=1
localhost:8000/api/account/balance?user_id=1&currency=USD
проверяем баланс аккауна id=2
localhost:8000/api/account/balance?user_id=2
выводим список транзакций аккауна id=1
localhost:8000/api/account/transacts?user_id=1
