package models

import (
	"gorm.io/gorm"
	"time"
)

type WalletOrder struct {
	gorm.Model
	TokenID     int
	Token       Token
	WalletID    int
	Wallet      Wallet
	Key         string
	DoneKey     string
	Type        string
	Currency    string
	SellAmount  float64
	BuyAmount   float64
	PaymentSum  float64
	ServiceCard string
	Revenue     float64
	IsConfirmed bool
	IsDone      bool
	CreatedAt   time.Time // column name is `created_at`
	UpdatedAt   time.Time // column name is `updated_at`
}

type WalletOrderRepository struct {
	DB *gorm.DB
}

func (r *WalletOrderRepository) Persist(walletOrder *WalletOrder) {
	r.DB.Save(walletOrder)
}

func (r *WalletOrderRepository) Delete(walletOrder *WalletOrder) {
	r.DB.Delete(walletOrder)
}

func (r *WalletOrderRepository) FindLastOrderByWalletId(walletId uint) *WalletOrder {
	var walletOrder WalletOrder
	res := r.DB.Last(&walletOrder, "wallet_id = ?", walletId)
	if res.Error != nil {
		return nil
	}
	return &walletOrder
}

func (r *WalletOrderRepository) FindOrderById(walletId uint) *WalletOrder {
	var walletOrder WalletOrder
	res := r.DB.First(&walletOrder, "id = ?", walletId)
	if res.Error != nil {
		return nil
	}
	return &walletOrder
}

func (r *WalletOrderRepository) FindLastOrderByTokenId(tokenId uint) *WalletOrder {
	var walletOrder WalletOrder
	res := r.DB.Last(&walletOrder, "tokenId = ?", tokenId)
	if res.Error != nil {
		return nil
	}
	return &walletOrder
}

func (r *WalletOrderRepository) FindOrderByKey(orderKey string) *WalletOrder {
	var walletOrder WalletOrder
	res := r.DB.Last(&walletOrder, "order_key = ?", orderKey)
	if res.Error != nil {
		return nil
	}
	return &walletOrder
}
