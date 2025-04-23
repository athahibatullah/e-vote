# Makefile for managing docker compose and Go commands

# Run the containers and build images
up:
	docker compose up -d --build

# Stop and remove containers
down:
	docker compose down

# Seed the database
seed:
	docker compose run --rm backend go run main.go seed

# View backend logs
logs:
	docker compose logs -f backend

# Rebuild only the backend service
rebuild-backend:
	docker compose build backend

# Restart backend service
restart-backend:
	docker compose restart backend
