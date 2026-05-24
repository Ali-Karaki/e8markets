#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

PROJECT=e8m-dev-test

copy_env_files() {
	for pair in \
		client/.env.example:client/.env.local \
		server/.env.example:server/.env.local \
		server/docker/db/.env.example:server/docker/db/.env.local
	do
		src="${pair%%:*}"
		dst="${pair#*:}"
		[ -f "$dst" ] || cp "$src" "$dst"
	done
}

stop() {
	docker rm -f e8m-client e8m-server e8m-postgres-db e8m-redis e8m-adminer 2>/dev/null || true
	(cd server/docker && docker compose down) 2>/dev/null || true
}

reset() {
	stop
	rm -rf server/docker/local-volumes
	mkdir -p server/docker/local-volumes/db-data server/docker/local-volumes/redis-data
}

dev() {
	copy_env_files

	cd server/docker
	docker compose up -d --wait postgres-db redis adminer
	cd "$ROOT"

	docker rm -f e8m-server 2>/dev/null || true
	docker build -t e8m-server server
	docker run -d --name e8m-server --network "$PROJECT" \
		--label com.docker.compose.project="$PROJECT" \
		--label com.docker.compose.service=server \
		--env-file server/.env.local \
		-p 8080:8080 \
		-v "$ROOT/server:/app" \
		-v e8m-go-mod-cache:/go/pkg/mod \
		e8m-server

	until curl -sf http://localhost:8080/health >/dev/null; do sleep 1; done

	trap stop EXIT INT TERM
	docker build -t e8m-client client
	echo ""
	echo "Ready"
	echo "  App  http://localhost:5173"
	echo "  API  http://localhost:8080/health"
	echo "  DB   http://localhost:9901/?pgsql=postgres-db&username=postgres&db=e8markets&ns=public"
	echo ""
	docker run --rm --name e8m-client \
		--label com.docker.compose.project="$PROJECT" \
		--label com.docker.compose.service=client \
		--env-file client/.env.local \
		-p 5173:5173 \
		-v "$ROOT/client:/app" \
		-v e8m-client-node-modules:/app/node_modules \
		e8m-client
}

case "${1:-dev}" in
	dev) dev ;;
	stop) stop ;;
	reset) reset ;;
	*)
		echo "usage: $0 [dev|stop|reset]" >&2
		exit 1
		;;
esac
