package controlle

import (
	"github.com/gofiber/fiber/v2"
	"go_shop/model"
	"go_shop/service/Cart"
	"go_shop/util"
)

type CartHandler struct {
	Cart *Cart.Cart
}

// AddItem 加入购物车主流程
func (ch *CartHandler) AddItem(c *fiber.Ctx) error {
	var ctm model.CartItem
	err := c.BodyParser(ctm)
	if err != nil {
		return util.Resp401(c, 401, "无效的参数请求")
	}
	if ctm.Number == 0 {
		return util.Resp401(c, 401, "加入购物车的商品数量为零")
	}
	err = ch.Cart.AddItem(ctm)
	if err != nil {
		return err
	}
	return util.Resp200(c, 200, "添加成功")
}
