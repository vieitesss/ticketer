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
up:
    docker compose up -d

# Stop all services
down:
    docker compose down

# View logs
logs service="":
    #!/usr/bin/env bash
    if [ -z "{{service}}" ]; then
        docker compose logs -f;
    else
        docker compose logs -f {{service}};
    fi

# Rebuild and restart all services
restart: down build-all up

# Clean up Docker resources
clean:
    docker compose down -v --rmi all --remove-orphans
