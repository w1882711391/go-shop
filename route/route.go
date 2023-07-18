package route

import (
	"github.com/gofiber/fiber/v2"
	"go_shop/controlle"
	"go_shop/service/Cart"
)

func RouterInit() *fiber.App {
	app := fiber.New()
	ch := &controlle.CartHandler{
		Cart: &Cart.Cart{},
	}
	cart := app.Group("/cart")
	{
		cart.Post("/addCart", ch.AddItem)
		cart.Post("addPd", ch.AddProduct)
	}
	return app
}
