package Product

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go_shop/dao"
	"go_shop/model"
	"time"
)

// AddProduct 添加商品核心逻辑
func AddProduct(pd *model.Product) error {
	if err := dao.DB.Table("products").Where("nick_name=? and user_id=?", pd.NickName, pd.UserId).Error; err != nil {
		return fmt.Errorf("cart.go：73 数据库中存在该商品：%v", err)
	}
	//开启一个事务 防止操作不一致
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := dao.DB.Create(&pd).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("cart.go：84 商品添加失败：%v", err)
	}
	// 提交事务

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交错误: %v", err)
	}
	logrus.Info("添加商品成功")
	return nil
}

// UpdateProduct 修改商品内容
func UpdateProduct(nickname string, userid string, pd model.Product) error {
	var product model.Product
	tx := dao.DB.Table("products").Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := dao.DB.Where("nike_name=? and user_id=?", nickname, userid).First(&product).Error; err != nil {
		return fmt.Errorf("查询不到此商品 : %v", err)
	}
	product.NickName = pd.NickName
	product.Stock = pd.Stock
	product.Price = pd.Price
	product.Description = pd.Description
	product.UpdatedAt = time.Now()

	if err := dao.DB.Save(&product).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("保存商品更改到数据库时出现错误: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交错误 %v", err)
	}
	logrus.Info("商品信息修改成功")
	return nil
}

// SearchProduct 查询商品
func SearchProduct(nickname string, userid string) (model.Product, error) {
	var product model.Product
	tx := dao.DB.Table("products").Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := dao.DB.Where("nike_name=? and user_id", nickname, userid).First(&product).Error; err != nil {
		return product, fmt.Errorf("查询不到此商品 : %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return product, fmt.Errorf("事务提交错误 %v", err)
	}
	//查到商品直接返回
	logrus.Info("查询商品成功")
	return product, nil
}

// DeleteProduct 删除商品 逻辑删除
func DeleteProduct(nickname, userid string) error {
	var product model.Product
	tx := dao.DB.Table("products").Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := dao.DB.Where("nike_name=? and user_id=?", nickname, userid).First(&product).Error; err != nil {
		return fmt.Errorf("查询不到此商品 : %v", err)
	}
	if err := dao.DB.Delete(&product); err != nil {
		tx.Rollback()
		return fmt.Errorf("删除失败: %v", err)
	}
	logrus.Info("删除成功")
	return nil
}
