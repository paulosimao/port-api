# Gens protobuffer
proto:
	protoc --go_out=./lib/proto --go_opt=paths=source_relative  --go-grpc_out=./lib/proto --go-grpc_opt=paths=source_relative  ./service.proto

bin2docker:
	GOOS=linux GOARCH=amd64 go build -o deploy/server/api cmd/server/main.go
	GOOS=linux GOARCH=amd64 go build -o deploy/cli/api cmd/cli/main.go

docker_build: bin2docker
	docker build -t port_server -f deploy/server/Dockerfile --platform linux/amd64 deploy/server
	docker build -t port_cli -f deploy/cli/Dockerfile --platform linux/amd64 deploy/cli
docker_net:
	- docker network create port_net
docker_run: bin2docker docker_stop docker_net docker_build
	docker run -d --net port_net --rm --name port_server --platform linux/amd64 -eADDR=:50052 port_server
	docker run -d --net port_net --rm --name port_cli  -p:8081:8080 -eADDR=:8080 -eGRPC_ADDR=port_server:50052 --platform linux/amd64 port_cli
docker_stop:
	- docker rm -f port_server
	- docker rm -f port_cli
docker_up:
	docker-compose -f ./deploy/docker-compose.yml up

test_up:
	curl -v -X POST -F "ports=@./ports.json" http://localhost:8081/

test_down:
	curl -v http://localhost:8081/

tests: test_up test_down

run_server:
	cd stage && ADDR=:50052 go run ../cmd/server/main.go

run_cli:
	cd stage && ADDR=:8080 GRPC_ADDR=localhost:50052 go run ../cmd/cli/main.go
