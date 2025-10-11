package http

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func NewServer() *fiber.App {
	app := fiber.New()

	if os.Getenv("IS_DEV") == "true" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: []string{"http://localhost:3000"},
			AllowHeaders: []string{"Content-Type", "Authorization"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			MaxAge:       3600,
		}))
	} else {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"https://kedada.fun"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowCredentials: true,
			MaxAge:           86400,
		}))
	}

	// app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${green}HTTP${reset} ${cyan}${method}${reset} ${green}${path}${reset} ${status} ${white}${latency}${reset}\n",
		TimeFormat: "15:04:05.000",
		TimeZone:   "Europe/Madrid",
		CustomTags: map[string]logger.LogFunc{
			"green": func(output logger.Buffer, c fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString("\x1b[92m") // Bright Green (ANSI 10)
			},
			"cyan": func(output logger.Buffer, c fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString("\x1b[96m") // Bright Cyan (ANSI 14)
			},
			"white": func(output logger.Buffer, c fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString("\x1b[97m") // Bright White (ANSI 15)
			},
			"reset": func(output logger.Buffer, c fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				return output.WriteString("\x1b[0m")
			},
		},
	}))

	return app
}
