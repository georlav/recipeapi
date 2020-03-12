run:
	GORACE="halt_on_error=1" go run -race cmd/api/main.go
build:
	go build -ldflags "-s -w" cmd/api/main.go
test:
	go test ./... -v -race -cover -count=1
import:
	go run cmd/import/main.go
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.23.6 golangci-lint run
lint-insecure:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.23.6 git config --global http.sslVerify false && golangci-lint run
db:
	docker-compose up -d &&	sleep 5
	docker exec mysql mysql -h 127.0.0.1 -u root -ppass -e "CREATE DATABASE IF NOT EXISTS recipes"
	docker exec mysql -h 127.0.0.1 -u root -ppass -e "CREATE DATABASE IF NOT EXISTS recipes_test"
	docker exec mysql -h 127.0.0.1 -u root -ppass recipes < recipes-schema.sql
	docker exec mysql -h 127.0.0.1 -u root -ppass recipes_test < recipes-schema.sql