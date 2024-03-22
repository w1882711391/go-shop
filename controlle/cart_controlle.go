package controlle

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_shop/model"
	"go_shop/service/Cart"
	"go_shop/util"
)

// AddItem 加入购物车主流程
func AddItem(c *fiber.Ctx) error {
	var ctm model.CartItem

	if err := c.BodyParser(&ctm); err != nil {
		return util.Resp401(c, "无效的参数请求")
	}

	if err := Cart.AddItem(c, ctm); err != nil {
		return util.Resp400(c, fmt.Sprintf("添加商品出现错误: %v", err))
	}
	return util.Resp200(c, "添加成功")
}

// IsOK 判断添加参数
func IsOK(pd model.Product) (bool, error) {
	if len(pd.Content) <= 0 || len(pd.NickName) <= 0 {
		return false, errors.New("商品名称和介绍需要大于0")
	}
	if pd.SalePrice <= 0 {
		return false, errors.New("商品价格需要大于0")
	}
	return true, nil
}

func DeleteItem(c *fiber.Ctx) error {
	nickname := c.FormValue("nick_name")
	userid := c.FormValue("user_id")

	if err := Cart.DeleteItem(nickname, userid); err != nil {
		return util.Resp400(c, fmt.Sprintf("删除商品出现错误: %v", err))
	}

	return util.Resp200(c, "删除成功")
}

func SearchItem(c *fiber.Ctx) error {
	userid := c.FormValue("user_id")

	items, err := Cart.SearchItem(userid)
	if err != nil {
		return util.Resp500(c, fmt.Errorf("查询错误 err:%v", err))
	}
	return util.Resp200(c, items, "查询成功")
}
