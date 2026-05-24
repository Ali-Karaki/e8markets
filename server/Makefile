build:
	go build -o ./cmd/main/tmp/main ./cmd/main/main.go

run:
	go run ./cmd/main/main.go

run-infra:
	cd docker && docker compose up

reset-infra:
	cd docker && docker compose down && rm -rf local-volumes && mkdir -p local-volumes/db-data local-volumes/redis-data
