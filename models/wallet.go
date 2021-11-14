package models

import (
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	gorm.Model
	TokenID     int
	Token       Token
	BalanceBTC  float64
	BalanceBNB  float64
	BalanceUSDT float64
	BalanceETH  float64
	CreatedAt   time.Time // column name is `created_at`
	UpdatedAt   time.Time // column name is `updated_at`
}

type WalletRepository struct {
	DB *gorm.DB
}

func (r *WalletRepository) Persist(wallet *Wallet) {
	r.DB.Save(wallet)
}

func (r *WalletRepository) Delete(wallet *Wallet) {
	r.DB.Delete(wallet)
}

func (r *WalletRepository) FindWalletByTokenId(tokenId uint) *Wallet {
	var wallet Wallet
	res := r.DB.First(&wallet, "token_id = ?", tokenId)
	if res.Error != nil {
		return nil
	}
	return &wallet
}

func (r *WalletRepository) FindById(id uint) *Wallet {
	var wallet Wallet
	res := r.DB.First(&wallet, "id = ?", id)
	if res.Error != nil {
		return nil
	}
	return &wallet
}
