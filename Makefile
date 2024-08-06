DATABASE_URL="postgresql://bompedido:bompedido@localhost:5432/bompedido?search_path=bompedido&sslmode=disable"

migrate:
	migrate -path=.sql/migrations -database $(DATABASE_URL) -verbose up

down:
	migrate -path=.sql/migrations -database $(DATABASE_URL) -verbose down

create-migration:
	migrate create -ext=sql -dir=.sql/migrations init

docker-build:
	docker build . --no-cache -t ederfmatos/bom-pedido-api:latest

docker-up:
	docker run --name bom-pedido-api --network host --env-file .env  ederfmatos/bom-pedido-api:latest