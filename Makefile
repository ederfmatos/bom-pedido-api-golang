docker-build:
	docker build . --no-cache -t ederfmatos/bom-pedido-api:latest

docker-up:
	docker run --name bom-pedido-api --network host --env-file .env  ederfmatos/bom-pedido-api:latest

lint:
	golangci-lint run