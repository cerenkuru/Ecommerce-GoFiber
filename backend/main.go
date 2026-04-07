package main

import (
	"log"
	"os"

	"github.com/cerenkuru/Ecommerce-GoFiber/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env dosyası yüklenemedi")
	}

	builder := bootstrap.NewAppBuilder()
	builder.WithDatabase()
	builder.WithDependencies()
	builder.WithMiddleware()
	builder.WithRoutes()

	app, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
