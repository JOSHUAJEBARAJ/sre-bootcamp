# Change these variables as necessary.
main_package_path = ./cmd/main.go
binary_name = server
.PHONY: help
include .env
## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	cd app && \
	go mod tidy -v && \
	go fmt ./...

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/${binary_name}

## test: run all tests
.PHONY: test
test:
	cd app && \
	go test -v -race -buildvcs ./...

## audit: run quality control checks


.PHONY: dev
dev:
	cd app && \
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "go build -o=/tmp/bin/${binary_name} ${main_package_path}" --build.bin "/tmp/bin/${binary_name}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


## Install: install the dependency
.PHONY: install 
install:
	cd app && \
	go mod tidy

.PHONY: start-db
start-db:
	docker run --rm --name postgres -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_DBNAME) -p 5432:5432 -v ./data:/var/lib/postgresql/data -d postgres:latest
.PHONY: stop-db
stop-db: 
	docker stop postgres

.PHONY: remove-db-data
remove-db-data: 
	rm -rf data

.PHONY: build
## build: build the application
build: install
	cd app && \
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

.PHONY: db-migrate
db-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
	migrate -path db/migrations -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@localhost:5432/$(DB_DBNAME)?sslmode=disable" up


.PHONY: docker-app-build
docker-app-build:
	cd app && \
	docker build -t student_api:$(IMAGE_TAG) .

.PHONY:docker-app-run
docker-app-run:
	@CONTAINER_HOST=$$(docker inspect postgres | jq -r '.[0].NetworkSettings.IPAddress'); \
	docker run  -d --rm --name app -p 8080:8080 -e DB_USERNAME=$(DB_USERNAME) -e DB_PASSWORD=$(DB_PASSWORD) -e DB_HOST=$$CONTAINER_HOST -e DB_PORT=$(DB_PORT) -e DB_DBNAME=$(DB_DBNAME) -e ENV=$(ENV) student_api:$(IMAGE_TAG)

.PHONY:docker-app-stop 
docker-app-stop:
	docker stop app 

.PHONY:docker-stop 
docker-stop: docker-app-stop stop-db

.PHONY:docker-compose-app-start 
docker-compose-app-start:
	docker-compose up -d db
	make db-migrate
	docker-compose up -d web 

.PHONY:docker-compose-app-stop 
docker-compose-app-stop:
	docker-compose down