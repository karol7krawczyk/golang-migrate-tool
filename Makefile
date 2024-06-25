.PHONY: build up down logs clean

export COMPOSE_PROJECT_NAME := migrate


up:
	@echo "Starting Docker containers..."
	docker compose up -d
	@make -s ps
	@make -s shell

build:
	@echo "Building Docker images..."
	docker compose build

down:
	@echo "Stopping Docker containers..."
	docker compose down --remove-orphans

logs:
	@echo "Fetching logs from Docker containers..."
	docker compose logs --tail=0 --follow

clean:
	@echo "Removing stopped containers and unused volumes..."
	docker compose down -v
	docker system prune -f

stop:
	docker compose stop

ps:
	docker compose ps

shell:
	docker compose exec -it golang bash

test:
	docker compose exec golang sh -c "go test"
