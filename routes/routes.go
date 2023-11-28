package routes

import (
	"github.com/alwilion/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/user")

	api.Post("/login", controllers.Login)
	api.Post("/register", controllers.Register)

	api.Get("/products", controllers.GetProductList)
	api.Get("/products/:id", controllers.GetProductById)
	api.Delete("/products/:id", controllers.DeleteProductById)
	api.Post("/products", controllers.AddProduct)
	api.Put("/products/:id", controllers.UpdateProduct)
}
