package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go_shop/controlle"

	"go_shop/util"
)

func Cors() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, x-token")
		c.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
func RouterInit() *fiber.App {
	app := fiber.New()
	app.Static("/image", "C:\\go-shop\\img\\")
	app.Use(Cors())
	app.Use(logger.New())
	// 购物车handle
	cart := app.Group("/cart", util.JWTMiddleware())
	{
		cart.Post("/cartItem", controlle.CartItem)
	}
	// 商品的handle
	app.Post("/image/upload", controlle.ImageUpload, util.JWTMiddleware())
	pd := app.Group("/product", util.JWTMiddleware())
	{
		pd.Get("/nolike/:pid", controlle.DeleteCart)
		pd.Get("/list", controlle.ViewProduct)
		pd.Get("/mylist", controlle.ViewCart)
		pd.Post("/addPd", controlle.AddProduct)
		pd.Get("/shelf/:id", controlle.DeleteProduct)
		pd.Post("/order/:pid", controlle.PlaceOrder)
		pd.Get("/viewOrder", controlle.ViewOrder)
		pd.Get("/myshop", controlle.ViewMyProduct)
	}
	// 用户handle
	user := app.Group("/user")
	{
		user.Post("/login", util.IsKick(), controlle.Login)
		user.Post("/updatePassword", util.JWTMiddleware(), controlle.UpdatePassword)
		user.Get("/myselfInformation", util.JWTMiddleware(), controlle.Information)
	}
	return app
}
