package controlle

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go_shop/dao"
	"go_shop/model"
	"go_shop/util"
	"gorm.io/gorm"
)

// CartItem 购物车流程
func CartItem(c *fiber.Ctx) error {
	type Req struct {
		ProductId uint `json:"product_id" form:"product_id"`
		IsLike    bool `json:"is_like" form:"is_like"`
	}
	var (
		req Req
	)
	req.ProductId = uint(c.QueryInt("product_id"))
	req.IsLike = c.QueryBool("is_like")

	userID, st := c.Locals("sid").(string)
	if !st {
		return errors.New("类型断言失败")
	}
	// 类型断言成功，可以将 strUserID 作为字符串类型使用

	if req.IsLike {
		ctm := model.CartItem{
			ProductId: req.ProductId,
			Sid:       userID,
			IsLike:    true,
		}

		err := dao.DB.Table("cart_items").Create(&ctm).Error
		if err != nil {
			return util.Resp400(c, "添加收藏错误")
		}

		dao.DB.Table("products").
			Where("id=?", req.ProductId).
			UpdateColumn("like_num", gorm.Expr("like_num+?", 1))

		return util.Resp200(c, "添加收藏成功")
	}
	// 不为true 取消收藏
	if err := dao.DB.Model(&model.CartItem{}).
		Where("product_id=? and sid=?", req.ProductId, userID).
		Delete(&model.CartItem{}).
		Error; err != nil {
		return util.Resp400(c, "取消收藏失败")
	}

	dao.DB.Table("products").
		Where("id=?", req.ProductId).
		UpdateColumn("like_num", gorm.Expr("like_num-?", 1))
	return util.Resp200(c, "取消收藏成功")
}

func DeleteCart(c *fiber.Ctx) error {
	pid := c.Params("pid")
	userID, st := c.Locals("sid").(string)
	if !st {
		return errors.New("类型断言失败")
	}
	// 类型断言成功，可以将 strUserID 作为字符串类型使用
	if err := dao.DB.Model(&model.CartItem{}).
		Where("product_id=? and sid=?", pid, userID).
		Delete(&model.CartItem{}).
		Error; err != nil {
		return util.Resp400(c, "取消收藏失败")
	}

	dao.DB.Table("products").
		Where("id=?", pid).
		UpdateColumn("like_num", gorm.Expr("like_num-?", 1))
	return util.Resp200(c, "取消收藏成功")
}

func ViewCart(c *fiber.Ctx) error {
	type Res struct {
		Pd    []model.Product `json:"product"`
		Total int             `json:"total"`
	}
	userID, st := c.Locals("sid").(string)
	if !st {
		return errors.New("类型断言失败")
	}
	// 类型断言成功，可以将 strUserID 作为字符串类型使用

	var (
		carts []model.CartItem
		pd    []model.Product
	)

	dao.DB.Model(&model.CartItem{}).
		Where("sid=? and deleted_at is NULL", userID).
		Find(&carts)
	pid := make([]uint, len(carts))
	for i, v := range carts {
		pid[i] = v.ProductId
	}
	dao.DB.Model(&model.Product{}).
		Where("id in ?", pid).
		Find(&pd)

	return util.Resp200(c, Res{
		pd,
		len(pd),
	})
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
