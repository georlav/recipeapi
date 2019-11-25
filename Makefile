run:
	GORACE="halt_on_error=1" go run -race cmd/api/main.go
build:
	go build -ldflags "-s -w" cmd/api/main.go
test:
	go test ./... -v -race -cover -count=1
import:
	go run cmd/import/main.go