package config

import (
	"errors"
	"log"
	"os"

	"eventhub/helper"
	"eventhub/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

const AppName = "EventHub"

type Config struct {
	AppPort string
	AppEnv  string

	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseURL string

	JWTSecret      string
	JWTExpireHours string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("config: file .env tidak ditemukan, memakai environment sistem")
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "eventhub"),
		DBPassword:  getEnv("DB_PASSWORD", "eventhub"),
		DBName:      getEnv("DB_NAME", "eventhub"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpireHours: getEnv("JWT_EXPIRE_HOURS", "24"),
	}
}

func NewApp(deps route.Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      AppName,
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	route.Setup(app, deps)

	app.Use(func(c *fiber.Ctx) error {
		return helper.NotFound("No endpoints found")
	})

	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	var appErr *helper.AppError
	if errors.As(err, &appErr) {
		resp := helper.WebResponse{Success: false, Message: appErr.Message}
		if appErr.Fields != nil {
			resp.Errors = appErr.Fields
		}
		return c.Status(appErr.Status).JSON(resp)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(helper.WebResponse{
		Success: false,
		Message: "Internal server error",
	})
}

func (c *Config) ServerAddr() string {
	return ":" + c.AppPort
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
