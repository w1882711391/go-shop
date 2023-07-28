package controlle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go_shop/model"
	"go_shop/service/Product"
	"go_shop/util"
)

// AddProduct 添加商品信息
func AddProduct(c *fiber.Ctx) error {
	var pd model.Product
	//绑定参数
	if err := c.BodyParser(&pd); err != nil {
		return util.Resp400(c, "无效的参数请求")
	}
	if ok, err := IsOK(pd); !ok {
		return util.Resp400(c, fmt.Sprintf("请求参数错误: %v", err))
	}
	if err := Product.AddProduct(&pd); err != nil {
		return util.Resp500(c, fmt.Sprintf("添加商品数据错误: %v", err))
	}
	return util.Resp200(c, 200, "添加成功")
}

// DeleteProduct 删除商品信息
func DeleteProduct(c *fiber.Ctx) error {
	nickname := c.FormValue("nike_name")
	if nickname == "" {
		return util.Resp401(c, "nickname为空")
	}
	if err := Product.DeleteProduct(nickname); err != nil {
		return util.Resp400(c, fmt.Errorf("删除失败 err:%v", err))
	}
	return util.Resp200(c, "删除成功")
}

// UpdateProduct 修改商品信息
func UpdateProduct(c *fiber.Ctx) error {
	nickname := c.FormValue("nike_name")
	var pd model.Product
	if err := c.BodyParser(&pd); err != nil {
		return util.Resp401(c, fmt.Errorf("没有绑定要修改的内容 err:%v", err))
	}
	if nickname == "" {
		return util.Resp401(c, "nickname为空")
	}

	if err := Product.UpdateProduct(nickname, pd); err != nil {
		return util.Resp400(c, fmt.Errorf("执行逻辑函数出错：%v", err))
	}
	return util.Resp200(c, "修改成功")
}

// SearchProduct 查询商品信息
func SearchProduct(c *fiber.Ctx) error {
	nickname := c.FormValue("nike_name")

	if nickname == "" {
		return util.Resp401(c, "nickname为空")
	}
	pd, err := Product.SearchProduct(nickname)
	if err != nil {
		return fmt.Errorf("查询错误: %v", err)
	}
	return util.Resp200(c, pd, "success")
}
