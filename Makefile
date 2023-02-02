setup:
	docker-compose up -d

setup-down:
	docker ps -a --format "{{.ID}} {{.Names}}" | grep cloudlens| awk '{print $$1}'| xargs docker stop | xargs docker rm -v

build:
	go build -o cloudlens main.go

run: build
	./cloudlens

populate: build
	./cloudlens lspop
