include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

migrate-force:
	@if [ -z "$(version)" ]; then echo "version is required (e.g. make migrate-force version=1)"; exit 1; fi
	migrate -path db/migrations -database "$(DB_URL)" force $(version)

migrate-create:
	@if [ -z "$(name)" ]; then echo "name is required (e.g. make migrate-create name=init)"; exit 1; fi
	migrate create -ext sql -dir db/migrations -seq $(name)

run:
	air

dev:
	air

build:
	cd cmd/api && go build -o ../../api.exe .
