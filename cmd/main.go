package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/kitae0522/gommunity/internal/controller"
	"github.com/kitae0522/gommunity/internal/model"
)

const port = ":8080"

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	dbconn := model.NewClient()
	if err := dbconn.Prisma.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func() {
		if err := dbconn.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	controller.EnrollRouter(app, dbconn)
	log.Fatal(app.Listen(port))
}
