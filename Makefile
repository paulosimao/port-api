# Gens protobuffer
proto:
	protoc --go_out=./lib/proto --go_opt=paths=source_relative  --go-grpc_out=./lib/proto --go-grpc_opt=paths=source_relative  ./service.proto

bin2docker:
	GOOS=linux GOARCH=amd64 go build -o deploy/server/api cmd/server/main.go
	GOOS=linux GOARCH=amd64 go build -o deploy/cli/api cmd/cli/main.go

docker_build:
	docker build -t server -f deploy/server/Dockerfile --platform linux/amd64 deploy/server
	docker build -t cli -f deploy/cli/Dockerfile --platform linux/amd64 deploy/cli

docker_run: bin2docker docker_build
	docker run -d --rm --name server --platform linux/amd64 server
	docker run -d --rm --name cli --platform linux/amd64 cli

docker_up:
	docker-compose -f ./deploy/docker-compose.yml up

test_up:
	curl -v -X POST -F "ports=@./ports.json" http://localhost:8080/

test_down:
	curl -v http://localhost:8080/

tests: test_up test_down

run_server:
	go run ./cmd/server/main.go

run_cli:
	go run ./cmd/cli/main.go
