.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./ ./...

.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows go build -o ./ ./...

.PHONY: run
run: build
	./todo

.PONY: mysql
mysql:
	docker-compose -f deployments/mysql.compose.yml up -d
