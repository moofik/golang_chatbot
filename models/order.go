package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Order struct {
	gorm.Model
	TokenID                   int
	Token                     Token
	MarketOrderKey            string
	MarketOrderType           string
	MarketOrderCurrency       string
	MarketOrderBuyAmount      float64
	MarketOrderBuyAddress     string
	MarketOrderSellAmount     float64
	MarketOrderSellCard       string
	MarketOrderPaymentSum     float64
	MarketOrderServiceAddress string
	MarketOrderServiceCard    string
	Revenue                   float64
	IsDone                    bool
	IsConfirmed               bool
	CreatedAt                 time.Time // column name is `created_at`
	UpdatedAt                 time.Time // column name is `updated_at`
}

type OrderRepository struct {
	DB *gorm.DB
}

func (r *OrderRepository) Persist(order *Order) {
	r.DB.Save(order)
}

func (r *OrderRepository) Delete(order *Order) {
	fmt.Println("DELETE ORDER WITH id " + strconv.Itoa(int(order.ID)))
	r.DB.Delete(order)
}

func (r *OrderRepository) DeleteByOrderKey(orderKey string) {
	r.DB.Delete(Order{}, "market_order_key = ? and is_confirmed IS NOT ?", orderKey, true)
}

func (r *OrderRepository) FindByOrderKey(orderKey string) *Order {
	var order Order
	res := r.DB.First(&order, "market_order_key = ? and deleted_at IS NULL", orderKey)
	if res.Error != nil {
		return nil
	}
	return &order
}