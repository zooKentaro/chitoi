run-watch:
	go run ./godo/main.go server --watch

docker-dev-start:
	docker-compose -f ./docker/dev/docker-compose.yml -p chitoi up -d

docker-dev-stop:
	docker-compose -f ./docker/dev/docker-compose.yml -p chitoi stop

import-masterdata:
	go run ./masterdata/hamster/main.go

migration-force:
	go run ./database/migration/main.go -deploy

generate-schema:
	go run ./database/ddlmaker/main.go
