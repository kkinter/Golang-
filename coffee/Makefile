DB_DOCKER_CONTAINER=coffee_db
BINARY_NAME=coffeapi
PORT = "8000"
DSN="host=localhost port=5432 user=root password=secret dbname=coffedb sslmode=disable timezone=UTC connect_timeout=5"
DB_DSN="postgres://root:secret@localhost:5432/coffeedb?sslmode=disable"

postgres:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.5-alpine

createdb:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=root --owner=root coffeedb

psql:
	docker exec -it coffee_db psql -U root

create_migrations:
	migrate create -seq -ext=.sql -dir=./migrations init

migrate-up:
	migrate -path=./migrations -database=${DB_DSN} up

migrate-down:
	migrate -path=./migrations -database=${DB_DSN} down

stop_containers:
	@echo "Stopping other docker containers"
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers.."; \
		docker stop $$(docker ps -q); \
	else \
		echo "no active containers found..."; \
	fi

start-docker:
	docker start ${DB_DOCKER_CONTAINER}

build:
	@echo "Building backend api binary"
	go build -o ${BINARY_NAME} cmd/server/*.go
	@echo "Binary built!"

run: build stop_containers start-docker
	@echo "Starting api"
	@env PORT=${PORT} DSN=${DSN} ./${BINARY_NAME} &
	@echo "api started!"

stop:
	@echo "Stopping backend"
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped backend"

start: run

restart: stop start
	

