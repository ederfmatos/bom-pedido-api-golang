package main

import (
	"bom-pedido-api/application/usecase/auth"
	"bom-pedido-api/infra/factory"
	"bom-pedido-api/infra/http"
	"bom-pedido-api/infra/registry"
	"bom-pedido-api/presentation/rest/handler"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func main() {
	_ = os.Setenv("GOOGLE_AUTH_URL", "https://www.googleapis.com/oauth2/v2/userinfo")

	database, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer database.Close()

	applicationFactory := factory.NewApplicationFactory(database)
	registry.RegisterDependency("GoogleAuthenticateCustomerUseCase", auth.NewGoogleAuthenticateCustomerUseCase(applicationFactory))

	server := http.NewHttpServer()
	server.HandleFunc("POST /auth/google/customer", handler.GoogleAuthCustomerHandler)
	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
