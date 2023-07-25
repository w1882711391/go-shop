package controlle

import (
	"errors"
	"fmt"
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

	if err := c.BodyParser(&ctm); err != nil {
		return util.Resp401(c, "无效的参数请求")
	}

	if ctm.Number <= 0 {
		return util.Resp400(c, "加入购物车的商品数量必须大于0")
	}

	if err := ch.Cart.AddItem(c, ctm); err != nil {
		return util.Resp400(c, fmt.Sprintf("添加商品出现错误: %v", err))
	}
	return util.Resp200(c, 200, "添加成功")
}

// AddProduct 添加商品信息
func (ch *CartHandler) AddProduct(c *fiber.Ctx) error {
	var pd model.Product
	//绑定参数
	if err := c.BodyParser(&pd); err != nil {
		return util.Resp400(c, "无效的参数请求")
	}
	if ok, err := ch.IsOK(pd); !ok {
		return util.Resp400(c, fmt.Sprintf("请求参数错误: %v", err))
	}
	if err := ch.Cart.AddProduct(&pd); err != nil {
		return util.Resp500(c, fmt.Sprintf("添加商品数据错误: %v", err))
	}
	return util.Resp200(c, 200, "添加成功")
}

// IsOK 判断添加参数
func (ch *CartHandler) IsOK(pd model.Product) (bool, error) {
	if len(pd.Description) <= 0 || len(pd.NickName) <= 0 {
		return false, errors.New("商品名称和介绍需要大于0")
	}
	if pd.Price <= 0 || pd.Stock <= 0 {
		return false, errors.New("商品价格和数量需要大于0")
	}
	return true, nil
}
