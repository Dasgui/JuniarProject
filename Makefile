all:
	$(MAKE) init-db
	$(MAKE) migrate
	go run ./cmd/
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