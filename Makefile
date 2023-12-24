
docker-up:
	docker compose -f docker-compose.yaml up --build

docker-down:
	docker compose down

docker-prune:
	docker system prune
