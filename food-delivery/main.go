package main

import (
	"flag"
	"food-delivery/database"
	"food-delivery/handlers"
	"food-delivery/messaging"
	"food-delivery/models"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
)

var (
	DSN   string
	PORT  string
	debug bool
	SEEDS string
)

func main() {
	service := "food-delivery"
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
	flag.Parse()

	DSN = os.Getenv("DSN")
	if DSN == "" {
		DSN = `host=localhost user=app password=app123 dbname=usersdb port=5432 sslmode=disable`
		log.Info().Msg(DSN)
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	SEEDS = os.Getenv("KAFKA_BROKERS")
	if SEEDS == "" {
		SEEDS = "localhost:19092,localhost:29092,localhost:39092"
	}

	db, err := database.GetConnection(DSN)

	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}
	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)

	msgUsersCreated := messaging.NewMessaging("orders.v1", strings.Split(SEEDS, ","))
	go msgUsersCreated.ProduceRecords()

	go msgUsersCreated.ConsumeRecords()

	app := fiber.New()

	prom := fiberprometheus.New(service)
	prom.RegisterAt(app, "/metrics") // exposes Prometheus metrics here
	app.Use(prom.Middleware)         // automatic request metrics

	app.Get("/", handlers.Root)
	app.Get("ping", handlers.Ping)
	app.Get("/health", handlers.Health)

	orderHandler := handlers.NewOrderHandler(database.NewOrderDB(db), database.NewOrderEventDB(db))
	order_group := app.Group("/api/v1/orders")
	order_group.Post("/", orderHandler.CreateOrder(msgUsersCreated))
	order_group.Get("/:order_id", orderHandler.GetOrder())

	//order_group.Get("/:id", orderHandler.GetOrderBy.GetUserBy)
	// user_group.Get("/all/:limit/:offset", userHandler.GetUsersByLimit)

	// order_group := app.Group("/api/v1/users/orders")
	// order_group.Post("/", userHandler.CreateOrder)

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.Orders{}, &models.OrdersEvent{})
}
