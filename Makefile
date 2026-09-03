all:
	docker-compose up --build
init-db:
	docker compose up -d

migrate:
	goose up

drop-bd:
	goose down

clean:
	$(MAKE) drop-bd
	docker compose down -v

test:
	go test -v -race ./internal/handler/product/...

swag:
	swag init -g cmd/main.go -o docs