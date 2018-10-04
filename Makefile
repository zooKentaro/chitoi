docker-dev-start:
	docker-compose -f ./docker/dev/docker-compose.yml -p chitoi up -d

docker-dev-stop:
	docker-compose -f ./docker/dev/docker-compose.yml -p chitoi stop
