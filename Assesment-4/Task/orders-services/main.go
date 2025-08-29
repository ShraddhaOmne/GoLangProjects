package main

import (
	"flag"
	"os"
	"users-service/database"
	"users-service/handlers"
	"users-service/messaging"
	"users-service/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	DSN   string
	PORT  string
	debug bool
	Seeds []string = []string{"localhost:19092", "localhost:29092", "localhost:39092"}
)

func main() {
	service := "users-service"
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
	flag.Parse()
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	DSN = os.Getenv("DSN")
	if DSN == "" {
		DSN = `host=localhost user=app password=app123 dbname=usersdb port=5432 sslmode=disable`
		log.Info().Msg(DSN)
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	db, err := database.GetConnection(DSN)

	if err != nil {
		//log.Fatal().Msg("unable to connect to the database..." + err.Error())
		log.Fatal().
			Err(err).
			Str("service", service).
			Msgf("unable to connect to the database %s", service)
	}
	log.Info().Str("service", service).Msg("database connection is established")
	Init(db)

	msgUsersCreated := messaging.NewMessaging("omnenext.users.created", Seeds)
	go msgUsersCreated.ProduceRecords()

	app := fiber.New()
	app.Get("/", handlers.Root)
	app.Get("ping", handlers.Ping)
	app.Get("/health", handlers.Health)

	userHandler := handlers.NewUserHandler(database.NewUserDB(db))
	user_group := app.Group("/api/v1/users")
	user_group.Post("/", userHandler.CreateUser(msgUsersCreated))
	user_group.Get("/:id", userHandler.GetUserBy)
	user_group.Get("/all/:limit/:offset", userHandler.GetUsersByLimit)

	order_group := app.Group("/api/v1/users/orders")
	order_group.Post("/", userHandler.CreateOrder)

	app.Listen(":" + PORT)

}

func Init(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Order{})
}
