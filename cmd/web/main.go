package main

import (
	"log"
	"os"

	"diwe-web/internal/handlers"
	"diwe-web/internal/pocketbase"
	"diwe-web/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New(fiber.Config{AppName: "DIWE B2B Web"})
	app.Static("/public", "./public")

	store := services.NewStoreWithSeed()
	pb := pocketbase.NewFromEnv()
	h := handlers.NewWebHandler(store, pb)
	h.Register(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
