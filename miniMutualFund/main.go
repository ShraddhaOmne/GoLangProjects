package main

import (
	"context"
	"flag"
	"miniMutualFund/database"
	"miniMutualFund/handlers"
	"miniMutualFund/messaging"
	"miniMutualFund/models"
	navredis "miniMutualFund/redis"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	service := "miniMutualFund"
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
	var ctx = context.Background()

	msgUsersCreated := messaging.NewMessaging("omnenest.mf.created", strings.Split(SEEDS, ","))
	go msgUsersCreated.ProduceRecords()

	go msgUsersCreated.ConsumeRecords()

	// initialize redis client
	rdb := navredis.NewRedisClient()
	handlers.StartNAVSimulator(ctx, rdb, []string{"SBI001", "ICICI001", "HDFC001", "SBI002", "ICICI002", "HDFC002", "SBI003", "ICICI003", "HDFC003"})

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	prom := fiberprometheus.New(service)
	prom.RegisterAt(app, "/metrics") // exposes Prometheus metrics here
	app.Use(prom.Middleware)         // automatic request metrics

	app.Get("/", handlers.Root)
	app.Get("ping", handlers.Ping)
	app.Get("/health", handlers.Health)

	orderHandler := handlers.NewOrderHandler(database.NewOrderDB(db))
	authHandler := &handlers.AuthHandler{}

	app.Post("/login", authHandler.Login)
	order_group := app.Group("/")
	order_group.Post("orders", orderHandler.CreateOrder(msgUsersCreated, rdb))
	order_group.Get("orders", orderHandler.GetAllOrders())

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.PlaceOrder{})
}
