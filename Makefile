start:
	docker-compose -f ./db/docker-compose.yml up -d && go build . && ./pagination

start_db:
	docker-compose -f ./db/docker-compose.yml up -d

stop:
	docker-compose -f ./db/docker-compose.yml down -v && rm -rf ./db/postgres-data

test:
	go test ./...