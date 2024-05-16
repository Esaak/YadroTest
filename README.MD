# Cтруктура проекта
````
YadroTest/
├── cmd/
│   ├── main.go
│   └── main_test.go
├── internal/
│   ├── model/
│   │   └── model.go
│   ├── process/
│   │   └── process.go
│   ├── root
│   │   └── root.go
│   └── usecase/
│       └── utils.go
├── tests/
│   └── test_data/
│       ├──test_base.go
│       └──test_to_end.go
├── Docker
├── README.md
├── go.mod
└── go.sum
````
## Запуск
``go run ./cmd/main.go <file_name>.txt``
## Запуск тестов
``go test -v ./...``
## Сборка docker
`` docker build -t <container-name> .``
## Запуск docker
`` docker run <container-name> <file_name>.txt`` 

P.S. Желательно добавлять новые файлы в папку: `tests/test_data/`