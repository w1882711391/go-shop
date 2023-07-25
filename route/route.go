package route

import (
	"github.com/gofiber/fiber/v2"
	"go_shop/controlle"
	"go_shop/util"
)

func RouterInit() *fiber.App {
	app := fiber.New()

	cart := app.Group("/cart", util.JWTMiddleware())
	{
		cart.Post("/addCart", controlle.AddItem)
		cart.Post("/addPd", controlle.AddProduct)
	}
	user := app.Group("/user")
	{
		user.Post("/create", controlle.Register)
		user.Post("/login", controlle.Login)
	}
	return app
}
