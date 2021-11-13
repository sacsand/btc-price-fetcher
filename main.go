package main

import (
	"log"
	"os"
	"test-server/config"
	"test-server/config/databases"
	"test-server/middleware"
	"test-server/pkg/coinprice"
	"test-server/routes"
	runners "test-server/runners"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := Setup()
	log.Fatal(app.Listen(":" + os.Getenv("HOST_PORT")))
}

// Setup - set middleware, router
func Setup() *fiber.App {
	config.InitConfig()  // load .env
	app := fiber.New()   // Create New Fiber APP
	middleware.Init(app) // Initialize middleware by passinf fiber app instance

	mySqlConnection := databases.InitMySql() // Create single Databse Connection
	api := app.Group("/api")                 // create api group

	/*
	   |--------------------------------------------------------------------------
	   | Repository and Service Registery
	   |--------------------------------------------------------------------------
	   |
	   |
	   |
	*/

	coinPriceRepository := coinprice.NewRepo(mySqlConnection)     // Init coinPrice Repository
	coinPriceService := coinprice.NewService(coinPriceRepository) // Init coinPrice Service

	/*
	   |--------------------------------------------------------------------------
	   | Routes Registery
	   |--------------------------------------------------------------------------
	   |
	   | Inject Service to the routes(Handlers)
	   |
	*/

	routes.CoinPriceRouter(api, coinPriceService)

	/*
	   |--------------------------------------------------------------------------
	   | JOBS Registery
	   |--------------------------------------------------------------------------
	   |
	   | Inject Service to the Jobs(Handlers)
	   |
	*/

	runners.CoinPriceRunner(coinPriceService) // Runner for fetching BTC data.

	// Default Error Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(401).SendString("Not Found " + c.BaseURL())
	})

	return app

}
