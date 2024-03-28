package controlle

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go_shop/dao"
	"go_shop/model"
	"go_shop/util"
	"math/rand"
	"strconv"
	"time"
)

const path = "./img/"

// AddProduct 添加商品信息
func AddProduct(c *fiber.Ctx) error {
	var pd model.Product
	userID, _ := c.Locals("sid").(string)
	//绑定参数
	if err := c.BodyParser(&pd); err != nil {
		return util.Resp400(c, "无效的参数请求")
	}
	if ok, err := IsOK(pd); !ok {
		return util.Resp400(c, fmt.Sprintf("请求参数错误: %v", err))
	}
	if err := dao.DB.Table("products").Where("nick_name=? and sid=?", pd.NickName, pd.Sid).Error; err != nil {
		return util.Resp400(c, "数据库钟存在该商品")
	}
	//开启一个事务 防止操作不一致
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	pd.LikeNum = 0
	var user model.User

	dao.DB.Table("users").
		Where("sid=?", userID).
		First(&user)
	pd.Sid = user.UserId
	pd.SName = user.UserName
	if err := tx.Create(&pd).Error; err != nil {
		tx.Rollback()
		return util.Resp400(c, "商品添加失败")
	}
	// 提交事务

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return util.Resp400(c, "商品添加失败")
	}
	logrus.Info("添加商品成功")
	return util.Resp200(c, 200, "添加成功")
}

// DeleteProduct 删除商品信息
func DeleteProduct(c *fiber.Ctx) error {

	pid := c.Params("id")

	userID, st := c.Locals("sid").(string)
	if !st {
		return util.Resp400(c, "类型断言失败")
	}
	tx := dao.DB.Table("products").Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := dao.DB.Where("id=? and sid=?", pid, userID).Delete(&model.Product{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除失败: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("数据库提交失败 ：%v", err)
	}

	return util.Resp200(c, "删除成功")
}

func ViewProduct(c *fiber.Ctx) error {

	type Req struct {
		Title    string `json:"title" form:"title"`
		PageNum  int64  `json:"pageNum" form:"pageNum"`
		PageSize int64  `json:"pageSize" form:"pageSize"`
	}
	type Res struct {
		Pd    []model.Product `json:"product"`
		Total int             `json:"total"`
	}
	var (
		pd  []model.Product
		req Req
		num int64
	)
	c.QueryParser(&req)
	tx := dao.DB.Table("products")
	tx2 := dao.DB.Model(&model.Product{})
	if req.Title != "" {
		tx.Where(fmt.Sprintf("nick_name like '%%%s%%'", req.Title))
		tx2.Where(fmt.Sprintf("nick_name like '%%%s%%'", req.Title))
	}
	err := tx.Limit(int(req.PageSize)).
		Offset(int((req.PageNum - 1) * req.PageSize)).
		Find(&pd).
		Error
	tx2.Count(&num)
	if err != nil {
		return util.Resp400(c, "商品数据库查询失败")
	}
	userID, _ := c.Locals("sid").(string)
	for i, v := range pd {
		var count int64 = 0
		dao.DB.Table("cart_items").
			Where("product_id=? and sid=?", v.ID, userID).
			Count(&count)
		if count != 0 {
			v.Like = true
		}
		var num int64
		dao.DB.Model(&model.CartItem{}).Where("product_id = ? and sid = ? and deleted_at is NULL", v.ID, userID).Count(&num)
		//fmt.Println(count)
		if num > 0 {
			pd[i].Like = true
		}

	}

	res := Res{
		Pd:    pd,
		Total: int(num),
	}
	return util.Resp200(c, res)
}

func PlaceOrder(c *fiber.Ctx) error {
	var (
		pd  model.Product
		err error
	)
	pid := c.Params("pid")
	userID, _ := c.Locals("sid").(string)
	dao.DB.Model(&model.Product{}).
		Where("id=?", pid).
		First(&pd)
	timestamp := time.Now().Format("20060102150405") // 当前时间，格式为年月日时分秒
	randomPart := generateRandomNumber(6)            // 生成长度为6的随机数部分

	if pd.Sid == userID {
		return util.Resp400(c, "禁止购买个人发布商品")
	}
	orderNumber := fmt.Sprintf("%d-%s-%s", pd.ID, timestamp, randomPart)
	var user model.User
	dao.DB.Model(&model.User{}).
		Where("sid=?", userID).
		First(&user)
	od := model.Order{
		OrderID:      orderNumber,
		Sid:          pd.Sid,
		SName:        pd.SName,
		ProductName:  pd.NickName,
		ProductID:    strconv.Itoa(int(pd.ID)),
		OrderStatus:  "已支付",
		OrderTime:    time.Now(),
		PaymentTime:  time.Now(),
		PayType:      "微信支付",
		ShippingTime: time.Now(),
		OId:          userID,
		OName:        user.UserName,
		Address:      "沈阳化工大学生活城E5-4128",
	}

	if err = dao.DB.Table("orders").
		Create(od).
		Error; err != nil {
		return util.Resp400(c, "订单获取失败")
	}

	if err := dao.DB.Model(&model.Product{}).
		Where("id=?", pid).
		UpdateColumn("is_order", true).
		Error; err != nil {
		return util.Resp400(c, "失败")
	}
	return util.Resp200(c, od)
}

func generateRandomNumber(length int) string {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func ImageUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return util.Resp400(c, "获取图片失败")
	}

	s := Rand() + ".png"
	err = c.SaveFile(file, path+s)
	if err != nil {
		return util.Resp400(c, err)
	}
	return util.Resp200(c, fmt.Sprintf("/image/%s", s))
}

func Rand() string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func ViewOrder(c *fiber.Ctx) error {
	type Res struct {
		Od    []model.Order
		Total int64
	}
	userID, _ := c.Locals("sid").(string)
	var (
		order []model.Order
	)
	if err := dao.DB.Model(&model.Order{}).
		Where("o_id=?", userID).
		Find(&order).Error; err != nil {
		return util.Resp400(c, "订单查询失败")
	}

	return util.Resp200(c, Res{
		Od:    order,
		Total: int64(len(order)),
	})
}

func ViewMyProduct(c *fiber.Ctx) error {
	userID, _ := c.Locals("sid").(string)

	var pd []model.Product

	dao.DB.Model(&model.Product{}).
		Where("sid=?", userID).
		Find(&pd)

	return util.Resp200(c, pd)
}
