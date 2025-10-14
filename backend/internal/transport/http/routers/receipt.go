package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/vieitesss/ticketer/internal/transport/http/handlers"
)

func NewReceiptRouter(server fiber.Router, handler handlers.ReceiptHandler) {
	receipt := server.Group("/receipts")

	receipt.Post("/upload", handler.UploadAndProcess)
	receipt.Get("/", handler.ListReceipts)
	receipt.Get("/:id", handler.GetReceipt)
	receipt.Delete("/:id", handler.DeleteReceipt)

	// Item routes
	server.Put("/items/:itemId", handler.UpdateItem)
}
