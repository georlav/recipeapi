run:
	GORACE="halt_on_error=1" go run -race cmd/api/main.go
test:
	go test ./... -v -race -cover -count=1
build:
	go build -ldflags "-s -w" cmd/api/main.go
swagger:
	swag init --dir cmd/api --output api/swagger
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run
lint-insecure:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.24.0 git config --global http.sslVerify false && golangci-lint run
db-dev:
	docker cp api/recipes-schema.sql mysql:/recipes-schema.sql
	docker cp api/recipes-data.sql mysql:/recipes-data.sql
	docker exec mysql /bin/sh -c 'mysql -h 127.0.0.1 -u root -ppass -e "DROP DATABASE IF EXISTS recipes;"'
	docker exec mysql /bin/sh -c 'mysql -h 127.0.0.1 -u root -ppass -e "CREATE DATABASE recipes"'
	docker exec mysql /bin/sh -c "mysql -h 127.0.0.1 -u root -ppass recipes < recipes-schema.sql"
	docker exec mysql /bin/sh -c "mysql -h 127.0.0.1 -u root -ppass recipes < recipes-data.sql"
db-test:
	docker cp api/recipes-schema.sql mysql:/recipes-schema.sql
	docker exec mysql /bin/sh -c 'mysql -h 127.0.0.1 -u root -ppass -e "DROP DATABASE IF EXISTS recipes_test"'
	docker exec mysql /bin/sh -c 'mysql -h 127.0.0.1 -u root -ppass -e "CREATE DATABASE recipes_test"'
	docker exec mysql /bin/sh -c "mysql -h 127.0.0.1 -u root -ppass recipes_test < recipes-schema.sql"
	docker exec mysql /bin/sh -c 'mysql -h 127.0.0.1 -u root -ppass -e "DROP DATABASE IF EXISTS recipes_handlers_test"'
	docker exec mysql /bin/sh -c 'mysql -h 127.0.0.1 -u root -ppass -e "CREATE DATABASE recipes_handlers_test"'
	docker exec mysql /bin/sh -c "mysql -h 127.0.0.1 -u root -ppass recipes_handlers_test < recipes-schema.sql"
db:
	docker-compose up -d &&	sleep 5
	make db-dev
	make db-test