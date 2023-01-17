setup:
	docker-compose up -d

setup-down:
	docker-compose down

build:
	go build -o cloudlens main.go

run: build
	./cloudlens start

populate: build
	./cloudlens lspop
