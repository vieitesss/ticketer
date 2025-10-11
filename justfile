alias bb := build-backend
alias bf := build-frontend
alias ba := build-all
alias rb := run-backend
alias rf := run-frontend

_default:
    just -l

# Build backend Docker image
build-backend:
    docker compose build backend

# Build frontend Docker image
build-frontend:
    docker compose build frontend

# Build all Docker images
build-all: build-backend build-frontend

# Run backend container
run-backend *cmd:
    docker compose run --rm backend {{cmd}}

# Run frontend container
run-frontend *cmd:
    docker compose run --rm frontend {{cmd}}

# Start all services with docker-compose
up *env:
	#!/usr/bin/env bash
	if [ "{{env}}" = "pro" ]; then
		docker compose -f docker-compose.yml -f docker-compose.pro.yml up -d;
		exit 0;
	fi

	docker compose up -d

# Stop all services
down *service:
	docker compose down {{service}}

# View logs
logs *service:
	docker compose logs {{service}} -f

# Rebuild and restart all services
restart:
	just down
	just build-all up

# Clean up Docker resources
clean:
	docker compose down -v --rmi all --remove-orphans
