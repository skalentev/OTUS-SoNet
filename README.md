# OTUS-SoNet

Прототип социальной сети для обучения по курсу Highload Architect

## Требования к ПО
- Go v1.18 or higher
- MySQL
- Docker
- Git

## Установка Dev
```bash
git clone https://github.com/skalentev/OTUS-SoNet
cd OTUS-SoNet
go mod tidy
sudo docker compose up -d
go run main.go
```

## Установка Prod
export GIN_MODE=release
```bash
git clone https://github.com/skalentev/OTUS-SoNet
cd OTUS-SoNet
go mod tidy
sudo docker compose up -d
export GIN_MODE=release
go run main.go
```

## Тестирование
Реализованы следующие запросы:


## License
Можно заимствовать без ограничений, о найденных ошибках и предложениях просьба сообщать на skalentev@gmail.com