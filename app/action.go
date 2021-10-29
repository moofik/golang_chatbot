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
	"net/http"
	"strconv"
	"time"
)

//app package actions

type ActionRegistry struct {
	DB *gorm.DB
}

func (ar *ActionRegistry) ActionRegistryHandler(name string, params map[string]string) runtime.Action {
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
		Model:                  gorm.Model{},
		TokenID:                int(t.GetId()),
		MarketOrderKey:         orderNumber,
		MarketOrderType:        extras["market_order_type"],
		MarketOrderCurrency:    extras["market_order_currency"],
		MarketOrderBuyAmount:   buyAmount,
		MarketOrderBuyAddress:  extras["market_order_buy_address"],
		MarketOrderPaymentSum:  float64(paymentSum),
		MarketOrderServiceCard: cardNumber,
		Revenue:                float64(paymentSum - actualPrice),
		IsDone:                 false,
		CreatedAt:              time.Time{},
		UpdatedAt:              time.Time{},
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
		Model:                     gorm.Model{},
		TokenID:                   int(t.GetId()),
		MarketOrderKey:            orderNumber,
		MarketOrderType:           extras["market_order_type"],
		MarketOrderCurrency:       extras["market_order_currency"],
		MarketOrderSellAmount:     sellAmount,
		MarketOrderSellCard:       extras["market_order_sell_card"],
		MarketOrderPaymentSum:     float64(paymentSum),
		MarketOrderServiceAddress: serviceAddress,
		Revenue:                   float64(actualPrice - paymentSum),
		IsDone:                    false,
		CreatedAt:                 time.Time{},
		UpdatedAt:                 time.Time{},
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

	if order.MarketOrderType == "Купить" {
		text = fmt.Sprintf(
			"Новый заказ %s: покупка %f %s за %d руб., адрес кошелька покупки - %s, карта на которую пользователь отправит деньги - %s",
			order.MarketOrderKey,
			order.MarketOrderBuyAmount,
			order.MarketOrderCurrency,
			int(order.MarketOrderPaymentSum),
			order.MarketOrderBuyAddress,
			order.MarketOrderServiceCard,
		)
	} else if order.MarketOrderType == "Продать" {
		text = fmt.Sprintf(
			"Новый заказ %s: продажа %f %s за %d руб., карта пользователя для вывода средств - %s, адрес кошелька куда пользователь отправит валюту - %s",
			order.MarketOrderKey,
			order.MarketOrderSellAmount,
			order.MarketOrderCurrency,
			int(order.MarketOrderPaymentSum),
			order.MarketOrderSellCard,
			order.MarketOrderServiceAddress,
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
