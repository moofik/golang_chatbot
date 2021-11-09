package app

import (
	"bot-daedalus/bot/command"
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

//app package actions

type ActionRegistry struct {
	DB *gorm.DB
}

func (ar *ActionRegistry) ActionRegistryHandler(name string, params map[string]string) runtime.Action {
	if name == "show_wallet" {
		return &ShowWallet{
			WalletRepository: &models.WalletRepository{DB: ar.DB},
		}
	}

	if name == "create_wallet" {
		return &CreateWallet{
			WalletRepository: &models.WalletRepository{DB: ar.DB},
		}
	}

	if name == "cancel_order_data" {
		return &CancelOrderData{
			OrderRepository: &models.OrderRepository{DB: ar.DB},
		}
	}

	if name == "confirm_market_order" {
		return &ConfirmMarketOrder{
			OrderRepository: &models.OrderRepository{DB: ar.DB},
		}
	}

	if name == "calculate_market_buy_order" {
		return &CalculateMarketBuyOrder{
			OrderRepository: &models.OrderRepository{DB: ar.DB},
		}
	}

	if name == "calculate_market_sell_order" {
		return &CalculateMarketSellOrder{
			OrderRepository: &models.OrderRepository{DB: ar.DB},
		}
	}

	if name == "calculate_wallet_buy_order" {
		return &CalculateWalletBuyOrder{
			WalletOrderRepository: &models.WalletOrderRepository{DB: ar.DB},
			WalletRepository:      &models.WalletRepository{DB: ar.DB},
		}
	}

	if name == "notify_wallet_buy_order" {
		return &NotifyWalletBuyOrder{
			WalletOrderRepository: &models.WalletOrderRepository{DB: ar.DB},
			WalletRepository:      &models.WalletRepository{DB: ar.DB},
		}
	}

	return nil
}

type CalculateMarketBuyOrder struct {
	OrderRepository *models.OrderRepository
}

func (a *CalculateMarketBuyOrder) GetName() string {
	return "calculate_market_buy_order"
}

func (a *CalculateMarketBuyOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	buyAmount, err := strconv.ParseFloat(extras["market_order_buy_amount"], 64)
	if err != nil {
		return nil
	}

	paymentSum, actualPrice, err := ConvertCrypto(extras["market_order_currency"], "RUB", buyAmount, true)
	if err != nil {
		return nil
	}

	orderNumber := uuid.New().String()
	cardNumber := GetCardNumber()

	order := &models.Order{
		Model:       gorm.Model{},
		TokenID:     int(t.GetId()),
		Key:         orderNumber,
		Type:        extras["market_order_type"],
		Currency:    extras["market_order_currency"],
		BuyAmount:   buyAmount,
		BuyAddress:  extras["market_order_buy_address"],
		PaymentSum:  float64(paymentSum),
		ServiceCard: cardNumber,
		Revenue:     float64(paymentSum - actualPrice),
		IsDone:      false,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	a.OrderRepository.Persist(order)
	extras["market_order_payment_sum"] = fmt.Sprintf("%d", paymentSum)
	extras["market_order_service_card"] = cardNumber
	extras["market_last_order_key"] = orderNumber
	t.SetExtras(extras)
	return nil
}

type CalculateMarketSellOrder struct {
	OrderRepository *models.OrderRepository
}

func (a *CalculateMarketSellOrder) GetName() string {
	return "calculate_marker_sell_order"
}

func (a *CalculateMarketSellOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	sellAmount, err := strconv.ParseFloat(extras["market_order_sell_amount"], 64)
	if err != nil {
		return nil
	}

	paymentSum, actualPrice, err := ConvertCrypto(extras["market_order_currency"], "RUB", sellAmount, false)
	if err != nil {
		return nil
	}

	orderNumber := uuid.New().String()
	serviceAddress := GetServiceAddress()

	order := &models.Order{
		Model:          gorm.Model{},
		TokenID:        int(t.GetId()),
		Key:            orderNumber,
		Type:           extras["market_order_type"],
		Currency:       extras["market_order_currency"],
		SellAmount:     sellAmount,
		SellCard:       extras["market_order_sell_card"],
		PaymentSum:     float64(paymentSum),
		ServiceAddress: serviceAddress,
		Revenue:        float64(actualPrice - paymentSum),
		IsDone:         false,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}

	a.OrderRepository.Persist(order)
	extras["market_order_payment_sum"] = fmt.Sprintf("%d", paymentSum)
	extras["market_order_service_address"] = serviceAddress
	extras["market_last_order_key"] = orderNumber
	t.SetExtras(extras)
	return nil
}

type ConfirmMarketOrder struct {
	OrderRepository *models.OrderRepository
}

func (a *ConfirmMarketOrder) GetName() string {
	return "confirm_market_order"
}

func (a *ConfirmMarketOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	orderKey := extras["market_last_order_key"]
	order := a.OrderRepository.FindByOrderKey(orderKey)
	order.IsConfirmed = true
	a.OrderRepository.Persist(order)

	text := ""

	if order.Type == "Купить" {
		text = fmt.Sprintf(
			"Новый заказ %s: покупка %f %s за %d руб., адрес кошелька покупки - %s, карта на которую пользователь отправит деньги - %s",
			order.Key,
			order.BuyAmount,
			order.Currency,
			int(order.PaymentSum),
			order.BuyAddress,
			order.ServiceCard,
		)
	} else if order.Type == "Продать" {
		text = fmt.Sprintf(
			"Новый заказ %s: продажа %f %s за %d руб., карта пользователя для вывода средств - %s, адрес кошелька куда пользователь отправит валюту - %s",
			order.Key,
			order.SellAmount,
			order.Currency,
			int(order.PaymentSum),
			order.SellCard,
			order.ServiceAddress,
		)
	} else {
		return nil
	}

	NotifyAdmins(text)
	return nil
}

type CancelOrderData struct {
	OrderRepository *models.OrderRepository
}

func (a *CancelOrderData) GetName() string {
	return "confirm_market_order"
}

func (a *CancelOrderData) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	key, keyOk := extras["market_last_order_key"]

	if keyOk {
		order := a.OrderRepository.FindByOrderKey(key)
		if !order.IsConfirmed {
			a.OrderRepository.Delete(order)
		}
	}

	extrasToDelete := []string{
		"market_order_type",
		"market_order_currency",
		"market_order_buy_amount",
		"market_order_buy_address",
		"market_order_sell_amount",
		"market_order_sell_card",
		"market_order_sell_amount",
		"market_order_payment_sum",
		"market_order_service_card,",
		"market_order_service_address",
		"market_last_order_key",
	}

	for _, extra := range extrasToDelete {
		_, ok := extras[extra]
		if ok {
			delete(extras, extra)
		}
	}

	t.SetExtras(extras)

	return nil
}

type ShowWallet struct {
	WalletRepository *models.WalletRepository
}

func (a *ShowWallet) GetName() string {
	return "show_wallet"
}

func (a *ShowWallet) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	text := "<b>- BTC</b>: {{.btc}} \n\n<b>- ETH</b>: {{.eth}} \n\n<b>- BNB</b>: {{.bnb}} \n\n<b>- USDT</b>: {{.usdt}}"
	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		return &runtime.GenericActionError{InnerError: err}
	}

	var tpl bytes.Buffer
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())
	data := map[string]float64{
		"btc":  wallet.BalanceBTC,
		"eth":  wallet.BalanceETH,
		"bnb":  wallet.BalanceBNB,
		"usdt": wallet.BalanceUSDT,
	}

	if err := tmpl.Execute(&tpl, data); err != nil {
		return &runtime.GenericActionError{InnerError: err}
	}

	result := tpl.String()
	err = p.SendTextMessage(result, runtime.ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	})

	if err != nil {
		return &runtime.GenericActionError{InnerError: err}
	}

	return nil
}

type CreateWallet struct {
	WalletRepository *models.WalletRepository
}

func (a *CreateWallet) GetName() string {
	return "show_wallet"
}

func (a *CreateWallet) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())

	if wallet == nil {
		wallet := &models.Wallet{
			Model:       gorm.Model{},
			TokenID:     int(t.GetId()),
			BalanceBTC:  0,
			BalanceBNB:  0,
			BalanceUSDT: 0,
			BalanceETH:  0,
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}
		a.WalletRepository.Persist(wallet)
	}

	return nil
}

type CalculateWalletBuyOrder struct {
	WalletOrderRepository *models.WalletOrderRepository
	WalletRepository      *models.WalletRepository
}

func (a *CalculateWalletBuyOrder) GetName() string {
	return "calculate_wallet_buy_order"
}

func (a *CalculateWalletBuyOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	buyAmount, err := strconv.ParseFloat(extras["wallet_order_buy_amount"], 64)
	if err != nil {
		return nil
	}

	paymentSum, actualPrice, err := ConvertCrypto(extras["wallet_order_currency"], "RUB", buyAmount, true)
	if err != nil {
		return nil
	}

	orderNumber := uuid.New().String()
	cardNumber := GetCardNumber()
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())

	order := &models.WalletOrder{
		Model:       gorm.Model{},
		TokenID:     int(t.GetId()),
		WalletID:    int(wallet.ID),
		Key:         orderNumber,
		Type:        extras["wallet_order_type"],
		Currency:    extras["wallet_order_currency"],
		BuyAmount:   buyAmount,
		PaymentSum:  float64(paymentSum),
		ServiceCard: cardNumber,
		Revenue:     float64(paymentSum - actualPrice),
		IsDone:      false,
		DoneKey:     uuid.NewString(),
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	a.WalletOrderRepository.Persist(order)
	extras["wallet_order_payment_sum"] = fmt.Sprintf("%d", paymentSum)
	extras["wallet_order_service_card"] = cardNumber
	extras["wallet_last_order_key"] = orderNumber
	t.SetExtras(extras)
	return nil
}

type NotifyWalletBuyOrder struct {
	WalletOrderRepository *models.WalletOrderRepository
	WalletRepository      *models.WalletRepository
}

func (a *NotifyWalletBuyOrder) GetName() string {
	return "notify_wallet_buy_order"
}

func (a *NotifyWalletBuyOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c command.Command,
) runtime.ActionError {
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())
	walletOrder := a.WalletOrderRepository.FindLastOrderByWalletId(wallet.ID)

	if walletOrder.Type == "Купить" {
		text := fmt.Sprintf(
			"Покупка на кошелек бота, новый заказ %s: покупка %f %s за %d руб., карта на которую пользователь отправит деньги - %s, код завершения сделки - %s",
			walletOrder.Key,
			walletOrder.BuyAmount,
			walletOrder.Currency,
			int(walletOrder.PaymentSum),
			walletOrder.ServiceCard,
			walletOrder.DoneKey,
		)

		NotifyAdmins(text)
	} else {
		return nil
	}

	return nil
}

func NotifyAdmins(text string) {
	url := "https://api.telegram.org/bot1799138792:AAGryx8c1D48yT8TAD5VCG1yzXs8k3tPtIc/sendMessage"

	reqBody := &runtime.TelegramOutgoingMessage{
		ChatID:    131231613,
		Text:      text,
		ParseMode: "HTML",
	}

	reqBytes, err := json.Marshal(reqBody)

	_, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(reqBytes),
	)

	if err != nil {
		panic(err)
	}
}
