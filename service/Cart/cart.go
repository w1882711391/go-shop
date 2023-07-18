package Cart

import "go_shop/model"

// Cart 购物车
type Cart struct {
	cart []model.CartItem
}

func (c *Cart) AddItem(p model.Product, num int16) error {
	NewProduct := model.CartItem{
		Product: p,
		Number:  num,
	}
	c.cart = append(c.cart, NewProduct)
}
