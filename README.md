# OTUS-SoNet

Прототип социальной сети для обучения по курсу Highload Architect

## Реализованы методы:
* /login
* /user/register
* /user/get/{id}
* /test/test
* /health
* /metrics

Используйте [POSTMAN коллекцию](https://github.com/skalentev/OTUS-SoNet/blob/main/OTUS-SoNet.postman_collection.json) для тестирования

## Требования к ПО для сборки и тестирования
- Go v1.18 or higher
- MySQL или Docker
- Git

## Запуск в Windows
1. Поднять БД MySQL (можно в докере, инструкция ниже)
2. Скачать репозиторий 
```bash
git clone https://github.com/skalentev/OTUS-SoNet
```
3. Перейти в папку с бинарником
   <code>cd OTUS-SoNet/bin</code>

4. в .env указать параметры подключекния к БД MySQL
5. Запуск
<code>./otus-sonet.exe</code>
6. Сервис поднимается на порту http://localhost:8080


## Запуск в Ubuntu
```bash
sudo apt install golang-go
git clone https://github.com/skalentev/OTUS-SoNet
cd OTUS-SoNet/
sudo docker compose up -d
go mod tidy
go test .
go run .
```




## Поднять MySQL в докере
```bash
git clone https://github.com/skalentev/OTUS-SoNet
cd OTUS-SoNet
sudo docker compose up -d
```

## Запуск в Go
```bash
git clone https://github.com/skalentev/OTUS-SoNet
cd OTUS-SoNet
go mod tidy
go test .
go run main.go
```

## Сборка приложения
```bash
git clone https://github.com/skalentev/OTUS-SoNet
cd OTUS-SoNet
go mod tidy
go test .
go build -o bin .
```

## Загрузка тестовых данных
В MySQL клиенте необходимо выполнить команды:
```sql
 SET GLOBAL local_infile=1;
 USE test1; 
 LOAD DATA LOCAL INFILE "/var/lib/mysql-files/people.csv" 
     INTO TABLE user 
     COLUMNS TERMINATED BY ',' 
     LINES TERMINATED BY '\n' 
     IGNORE 1 LINES 
     (second_name, first_name, @age, city) 
     SET id = UUID(), birthdate = date_sub(current_date(),INTERVAL @age YEAR);
```
