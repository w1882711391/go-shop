package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go_shop/controlle"
	"go_shop/util"
)

func RouterInit() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	cart := app.Group("/cart", util.JWTMiddleware())
	{
		cart.Post("/addCart", controlle.AddItem)
	}
	pd := app.Group("/product")
	{
		pd.Post("/addPd", controlle.AddProduct)
		pd.Post("/deletePd", controlle.DeleteProduct)
		pd.Post("/updatePd", controlle.UpdateProduct)
		pd.Get("/searchPd", controlle.SearchProduct)
	}
	user := app.Group("/user")
	{
		user.Post("/create", controlle.Register)
		user.Post("/login", controlle.Login)
	}
	return app
}
