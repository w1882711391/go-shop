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
	// 购物车handle
	cart := app.Group("/cart", util.JWTMiddleware())
	{
		cart.Post("/addCart", controlle.AddItem)
		cart.Post("/updateCart", controlle.UpdateItem)
		cart.Post("/searchCart", controlle.SearchItem)
		cart.Post("/deleteCart", controlle.DeleteItem)
	}
	// 商品的handle
	pd := app.Group("/product", util.JWTMiddleware(), util.IsMerchant())
	{
		pd.Post("/addPd", controlle.AddProduct)
		pd.Post("/deletePd", controlle.DeleteProduct)
		pd.Post("/updatePd", controlle.UpdateProduct)
		pd.Get("/searchPd", controlle.SearchProduct)
	}
	// 用户handle
	user := app.Group("/user")
	{
		user.Post("/kickUser", controlle.KickUser)
		user.Post("/create", controlle.Register)
		user.Post("/login", util.IsKick(), controlle.Login)
	}
	return app
}
