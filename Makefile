build:
	go build -o cloudlens main.go

populate: build
	./cloudlens lspop