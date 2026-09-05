run:
	docker compose up --build

docker-rebuild:
	docker compose down -v && \
	docker compose up --build

migrate:
	goose up

drop-bd:
	goose down

test:
	go generate ./internal/service/product/...
	go test -v -race ./internal/handler/product/...

swag:
	swag init -g cmd/main.go -o docs