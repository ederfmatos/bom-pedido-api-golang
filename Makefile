DATABASE_URL="postgresql://bompedido:bompedido@localhost:5432/bompedido?sslmode=disable&dbname=bompedido"

migrate:
	migrate -path=.resources/sql/migrations -database $(DATABASE_URL) -verbose up

down:
	migrate -path=.resources/sql/migrations -database $(DATABASE_URL) -verbose down

create-migration:
	migrate create -ext=sql -dir=.resources/sql/migrations init

docker-build:
	docker build . --no-cache -t ederfmatos/bom-pedido-api:latest

docker-up:
	docker run --name bom-pedido-api --network host --env-file .env  ederfmatos/bom-pedido-api:latest

lint:
	golangci-lint run