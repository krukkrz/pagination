start:
	docker-compose -f ./db/docker-compose.yml up -d

stop:
	docker-compose -f ./db/docker-compose.yml down -v && rm -rf ./db/postgres

test:
	go test ./...