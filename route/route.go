package route

import (
	"github.com/gofiber/fiber/v2"
	"go_shop/controlle"
	"go_shop/service/Cart"
	"go_shop/service/User"
	"go_shop/util"
)

func RouterInit() *fiber.App {
	app := fiber.New()

	ch := &controlle.CartHandler{
		Cart: &Cart.Cart{},
	}
	us := &controlle.UserHandler{
		User: &User.NewUser{},
	}
	cart := app.Group("/cart", util.JWTMiddleware())
	{
		cart.Post("/addCart", ch.AddItem)
		cart.Post("/addPd", ch.AddProduct)
	}
	user := app.Group("/user")
	{
		user.Post("/create", us.Register)
		user.Post("/login", us.Login)
	}
	return app
}
