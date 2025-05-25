docker-up:
	docker compose up

docker-down-v:
	docker compose down -v

logs:
	docker compose logs -f app

docker-stop:
	docker compose stop

docker-start:
	docker compose start

test:
	go test -v ./...

doc-up-n-c:
	docker compose build --no-cache && docker compose up

swag:
	swag init -g ./cmd/main.go --parseDependency --parseInternal
