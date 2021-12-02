package app

import (
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//app package actions

type ActionRegistry struct {
	DB *gorm.DB
}

func (ar *ActionRegistry) ActionRegistryHandler(name string, params map[string]interface{}) runtime.Action {
	if name == "show_wallet" {
		return &ShowWallet{
			Params:           params,
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
			OrderRepository:       &models.OrderRepository{DB: ar.DB},
			WalletOrderRepository: &models.WalletOrderRepository{DB: ar.DB},
		}
	}

	if name == "confirm_market_order" {
		return &ConfirmMarketOrder{
			OrderRepository:    &models.OrderRepository{DB: ar.DB},
			SettingsRepository: &models.SettingsRepository{DB: ar.DB},
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

	if name == "notify_wallet_order" {
		return &NotifyWalletOrder{
			WalletOrderRepository: &models.WalletOrderRepository{DB: ar.DB},
			WalletRepository:      &models.WalletRepository{DB: ar.DB},
			SettingsRepository:    &models.SettingsRepository{DB: ar.DB},
		}
	}

	if name == "calculate_wallet_sell_order" {
		return &CalculateWalletSellOrder{
			WalletOrderRepository: &models.WalletOrderRepository{DB: ar.DB},
			WalletRepository:      &models.WalletRepository{DB: ar.DB},
		}
	}

	if name == "calculate_wallet_exchange_order" {
		return &CalculateWalletExchangeOrder{
			WalletOrderRepository: &models.WalletOrderRepository{DB: ar.DB},
			WalletRepository:      &models.WalletRepository{DB: ar.DB},
		}
	}

	if name == "confirm_order" {
		return &ConfirmOrder{
			WalletRepository:      &models.WalletRepository{DB: ar.DB},
			WalletOrderRepository: &models.WalletOrderRepository{DB: ar.DB},
			OrderRepository:       &models.OrderRepository{DB: ar.DB},
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
	c runtime.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	buyAmountRaw := strings.Replace(extras["market_order_buy_amount"], ",", ".", -1)
	buyAmount, err := strconv.ParseFloat(buyAmountRaw, 64)
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
		DoneKey:     "market_" + uuid.NewString(),
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

func (a *CalculateMarketBuyOrder) GetAlias() string {
	return "a1"
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
	c runtime.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	amountRaw := strings.Replace(extras["market_order_sell_amount"], ",", ".", -1)
	sellAmount, err := strconv.ParseFloat(amountRaw, 64)
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
		DoneKey:        "market_" + uuid.NewString(),
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
	OrderRepository    *models.OrderRepository
	SettingsRepository *models.SettingsRepository
}

func (a *ConfirmMarketOrder) GetName() string {
	return "confirm_market_order"
}

func (a *ConfirmMarketOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c runtime.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	orderKey := extras["market_last_order_key"]
	order := a.OrderRepository.FindByOrderKey(orderKey)
	order.IsConfirmed = true
	a.OrderRepository.Persist(order)

	text := ""

	if order.Type == "Купить" {
		text = fmt.Sprintf(
			"Маркет: новый заказ %s, покупка %f %s за %d руб., адрес кошелька покупки - %s, карта на которую пользователь отправит деньги - %s, код завершения сделки - %s",
			order.Key,
			order.BuyAmount,
			order.Currency,
			int(order.PaymentSum),
			order.BuyAddress,
			order.ServiceCard,
			order.DoneKey,
		)
	} else if order.Type == "Продать" {
		text = fmt.Sprintf(
			"Маркет: новый заказ %s, продажа %f %s за %d руб., карта пользователя для вывода средств - %s, адрес кошелька куда пользователь отправит валюту - %s, код завершения сделки - %s",
			order.Key,
			order.SellAmount,
			order.Currency,
			int(order.PaymentSum),
			order.SellCard,
			order.ServiceAddress,
			order.DoneKey,
		)
	} else {
		return nil
	}

	settings := a.SettingsRepository.FindByScenarioName(p.GetScenarioName())

	if settings != nil {
		for _, id := range settings.GetTelegramAdminsIds() {
			for _, botToken := range settings.GetTelegramNotificationChannelsTokens() {
				NotifyAdmins(text, id, botToken)
			}
		}
	}

	return nil
}

type CancelOrderData struct {
	OrderRepository       *models.OrderRepository
	WalletOrderRepository *models.WalletOrderRepository
}

func (a *CancelOrderData) GetName() string {
	return "confirm_market_order"
}

func (a *CancelOrderData) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c runtime.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	key, keyOk := extras["market_last_order_key"]

	if keyOk {
		order := a.OrderRepository.FindByOrderKey(key)
		if order != nil && !order.IsConfirmed {
			a.OrderRepository.Delete(order)
		}
	}

	key, keyOk = extras["wallet_last_order_key"]

	if keyOk {
		walletOrder := a.WalletOrderRepository.FindOrderByKey(key)
		if walletOrder != nil && !walletOrder.IsConfirmed {
			a.WalletOrderRepository.Delete(walletOrder)
		}
	}

	extrasToDelete := []string{
		"market_order_type",
		"market_order_currency",
		"market_order_buy_amount",
		"market_order_buy_address",
		"market_order_sell_amount",
		"market_order_sell_card",
		"market_order_payment_sum",
		"market_order_service_card,",
		"market_order_service_address",
		"market_last_order_key",
		"wallet_last_order_key",
		"wallet_order_buy_amount",
		"wallet_order_currency",
		"wallet_order_payment_sum",
		"wallet_order_service_card",
		"wallet_order_client_card",
		"wallet_order_exchange_address",
		"wallet_order_exchange_amount",
		"wallet_order_sell_amount",
		"wallet_order_type",
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
	Params           map[string]interface{}
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
	c runtime.Command,
) runtime.ActionError {
	text := "<b>- BTC</b>: {{.btc}}  ~  {{.btc_rub}} руб. \n\n<b>- ETH</b>: {{.eth}}  ~  {{.eth_rub}} руб. \n\n<b>- BNB</b>: {{.bnb}}  ~  {{.bnb_rub}} руб. \n\n<b>- USDT</b>: {{.usdt}}  ~  {{.usdt_rub}} руб."
	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		return &runtime.GenericActionError{InnerError: err}
	}

	var tpl bytes.Buffer
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())
	_, BTCinRUB, _ := ConvertCrypto("BTC", "RUB", wallet.BalanceBTC, false)
	_, ETHinRUB, _ := ConvertCrypto("ETH", "RUB", wallet.BalanceETH, false)
	_, BNBinRUB, _ := ConvertCrypto("BNB", "RUB", wallet.BalanceBNB, false)
	_, USDTinRUB, _ := ConvertCrypto("USDT", "RUB", wallet.BalanceUSDT, false)

	data := map[string]string{
		"btc":      fmt.Sprintf("%f", wallet.BalanceBTC),
		"eth":      fmt.Sprintf("%f", wallet.BalanceETH),
		"bnb":      fmt.Sprintf("%f", wallet.BalanceBNB),
		"usdt":     fmt.Sprintf("%f", wallet.BalanceUSDT),
		"btc_rub":  fmt.Sprintf("%d", BTCinRUB),
		"eth_rub":  fmt.Sprintf("%d", ETHinRUB),
		"bnb_rub":  fmt.Sprintf("%d", BNBinRUB),
		"usdt_rub": fmt.Sprintf("%d", USDTinRUB),
	}

	if err := tmpl.Execute(&tpl, data); err != nil {
		return &runtime.GenericActionError{InnerError: err}
	}

	result := tpl.String()
	lastBotMessageId := uint(t.GetLastBotMessageId())
	err = p.SendTextMessage(result, runtime.ProviderContext{
		State:   s,
		Command: c,
		Token:   t,
	})

	if err != nil {
		return &runtime.GenericActionError{InnerError: err}
	}

	if a.Params != nil {
		if clear, ok := a.Params["clear_previous"]; ok && clear.(bool) && t.GetIsLastBotMessageRemovable() {
			runtime.DeleteMessage(t.GetChatId(), lastBotMessageId, p.GetConfig().Token)
		}

		if removable, ok := a.Params["removable"]; ok && !removable.(bool) {
			t.SetIsLastBotMessageRemovable(false)
		} else {
			t.SetIsLastBotMessageRemovable(true)
		}
	} else {
		t.SetIsLastBotMessageRemovable(true)
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
	c runtime.Command,
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
	c runtime.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	buyAmount, err := strconv.ParseFloat(
		strings.Replace(extras["wallet_order_buy_amount"], ",", ".", -1),
		64,
	)

	if err != nil {
		fmt.Println("error: " + err.Error())
		return &runtime.GenericActionError{
			InnerError: errors.Errorf("Формат валюты должен быть следующим ( с точкой разделителем ): 23.940"),
		}
	}

	paymentSum, actualPrice, err := ConvertCrypto(extras["wallet_order_currency"], "RUB", buyAmount, true)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return &runtime.GenericActionError{InnerError: err}
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
		DoneKey:     "wallet_" + uuid.NewString(),
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

type CalculateWalletSellOrder struct {
	WalletOrderRepository *models.WalletOrderRepository
	WalletRepository      *models.WalletRepository
}

func (a *CalculateWalletSellOrder) GetName() string {
	return "calculate_wallet_sell_order"
}

func (a *CalculateWalletSellOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c runtime.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	amountRaw := strings.Replace(extras["wallet_order_sell_amount"], ",", ".", -1)
	sellAmount, err := strconv.ParseFloat(amountRaw, 64)

	if err != nil {
		fmt.Println("error: " + err.Error())
		return &runtime.GenericActionError{InnerError: err}
	}

	paymentSum, actualPrice, err := ConvertCrypto(extras["wallet_order_currency"], "RUB", sellAmount, false)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return &runtime.GenericActionError{InnerError: err}
	}

	orderNumber := uuid.New().String()
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())

	order := &models.WalletOrder{
		Model:      gorm.Model{},
		TokenID:    int(t.GetId()),
		WalletID:   int(wallet.ID),
		Key:        orderNumber,
		Type:       extras["wallet_order_type"],
		Currency:   extras["wallet_order_currency"],
		SellAmount: sellAmount,
		PaymentSum: float64(paymentSum),
		ClientCard: extras["wallet_order_client_card"],
		Revenue:    float64(actualPrice - paymentSum),
		IsDone:     false,
		DoneKey:    "wallet_" + uuid.NewString(),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	a.WalletOrderRepository.Persist(order)
	extras["wallet_order_payment_sum"] = fmt.Sprintf("%d", paymentSum)
	extras["wallet_last_order_key"] = orderNumber
	t.SetExtras(extras)

	return nil
}

type CalculateWalletExchangeOrder struct {
	WalletOrderRepository *models.WalletOrderRepository
	WalletRepository      *models.WalletRepository
}

func (a *CalculateWalletExchangeOrder) GetName() string {
	return "calculate_wallet_exchange_order"
}

func (a *CalculateWalletExchangeOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c runtime.Command,
) runtime.ActionError {
	extras := t.GetExtras()
	amountRaw := strings.Replace(extras["wallet_order_exchange_amount"], ",", ".", -1)
	exchangeAmount, err := strconv.ParseFloat(amountRaw, 64)

	if err != nil {
		fmt.Println("error: " + err.Error())
		return &runtime.GenericActionError{InnerError: err}
	}

	if err != nil {
		fmt.Println("error: " + err.Error())
		return &runtime.GenericActionError{InnerError: err}
	}

	orderNumber := uuid.New().String()
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())

	order := &models.WalletOrder{
		Model:           gorm.Model{},
		TokenID:         int(t.GetId()),
		WalletID:        int(wallet.ID),
		Key:             orderNumber,
		Type:            extras["wallet_order_type"],
		Currency:        extras["wallet_order_currency"],
		ExchangeAmount:  exchangeAmount,
		ExchangeAddress: extras["wallet_order_exchange_address"],
		Revenue:         0,
		IsDone:          false,
		DoneKey:         "wallet_" + uuid.NewString(),
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	a.WalletOrderRepository.Persist(order)
	extras["wallet_last_order_key"] = orderNumber
	t.SetExtras(extras)

	return nil
}

type ConfirmOrder struct {
	WalletRepository      *models.WalletRepository
	WalletOrderRepository *models.WalletOrderRepository
	OrderRepository       *models.OrderRepository
}

func (a *ConfirmOrder) GetName() string {
	return "confirm_order"
}

func (a *ConfirmOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c runtime.Command,
) runtime.ActionError {
	input := c.GetInput()

	if input[0] == 'm' { // market_ prefix
		order := a.OrderRepository.FindByDoneKey(c.GetInput())
		order.IsDone = true
		a.OrderRepository.Persist(order)
	} else { // wallet_ prefix
		walletOrder := a.WalletOrderRepository.FindByDoneKey(c.GetInput())

		if walletOrder == nil {
			panic("WALLET ORDER NOT FOUND!!!")
		}

		walletOrder.IsDone = true
		a.WalletOrderRepository.Persist(walletOrder)

		wallet := a.WalletRepository.FindById(uint(walletOrder.WalletID))

		if walletOrder.Type == "Купить" {
			if walletOrder.Currency == "BTC" {
				wallet.BalanceBTC += walletOrder.BuyAmount
			} else if walletOrder.Currency == "ETH" {
				wallet.BalanceETH += walletOrder.BuyAmount
			} else if walletOrder.Currency == "BNB" {
				wallet.BalanceBNB += walletOrder.BuyAmount
			} else if walletOrder.Currency == "USDT" {
				wallet.BalanceUSDT += walletOrder.BuyAmount
			}
		} else if walletOrder.Type == "Продать" {
			if walletOrder.Currency == "BTC" {
				wallet.BalanceBTC -= walletOrder.SellAmount
			} else if walletOrder.Currency == "ETH" {
				wallet.BalanceETH -= walletOrder.SellAmount
			} else if walletOrder.Currency == "BNB" {
				wallet.BalanceBNB -= walletOrder.SellAmount
			} else if walletOrder.Currency == "USDT" {
				wallet.BalanceUSDT -= walletOrder.SellAmount
			}
		} else {
			if walletOrder.Currency == "BTC" {
				wallet.BalanceBTC -= walletOrder.ExchangeAmount
			} else if walletOrder.Currency == "ETH" {
				wallet.BalanceETH -= walletOrder.ExchangeAmount
			} else if walletOrder.Currency == "BNB" {
				wallet.BalanceBNB -= walletOrder.ExchangeAmount
			} else if walletOrder.Currency == "USDT" {
				wallet.BalanceUSDT -= walletOrder.ExchangeAmount
			}
		}

		a.WalletRepository.Persist(wallet)
	}

	return nil
}

type NotifyWalletOrder struct {
	WalletOrderRepository *models.WalletOrderRepository
	WalletRepository      *models.WalletRepository
	SettingsRepository    *models.SettingsRepository
}

func (a *NotifyWalletOrder) GetName() string {
	return "notify_wallet_order"
}

func (a *NotifyWalletOrder) Run(
	p runtime.ChatProvider,
	t runtime.TokenProxy,
	s *runtime.State,
	prev *runtime.State,
	c runtime.Command,
) runtime.ActionError {
	wallet := a.WalletRepository.FindWalletByTokenId(t.GetId())
	walletOrder := a.WalletOrderRepository.FindLastOrderByWalletId(wallet.ID)
	walletOrder.IsConfirmed = true
	a.WalletOrderRepository.Persist(walletOrder)
	settings := a.SettingsRepository.FindByScenarioName(p.GetScenarioName())

	if walletOrder.Type == "Купить" {
		text := fmt.Sprintf(
			"Кошелек: новый заказ %s: покупка %f %s за %d руб., карта на которую пользователь отправит деньги - %s, код завершения сделки - %s",
			walletOrder.Key,
			walletOrder.BuyAmount,
			walletOrder.Currency,
			int(walletOrder.PaymentSum),
			walletOrder.ServiceCard,
			walletOrder.DoneKey,
		)

		if settings != nil {
			for _, id := range settings.GetTelegramAdminsIds() {
				for _, botToken := range settings.GetTelegramNotificationChannelsTokens() {
					NotifyAdmins(text, id, botToken)
				}
			}
		}
	} else if walletOrder.Type == "Продать" {
		text := fmt.Sprintf(
			"Кошелёк: новый заказ %s: продажа %f %s за %d руб., карта на которую пользователь получит деньги - %s, код завершения сделки - %s",
			walletOrder.Key,
			walletOrder.SellAmount,
			walletOrder.Currency,
			int(walletOrder.PaymentSum),
			walletOrder.ClientCard,
			walletOrder.DoneKey,
		)

		if settings != nil {
			for _, id := range settings.GetTelegramAdminsIds() {
				for _, botToken := range settings.GetTelegramNotificationChannelsTokens() {
					NotifyAdmins(text, id, botToken)
				}
			}
		}
	} else if walletOrder.Type == "Перевод" {
		text := fmt.Sprintf(
			"Кошелёк: новый заказ %s: перевод %f %s, адрес кошелька перевода - %s, код завершения сделки - %s",
			walletOrder.Key,
			walletOrder.ExchangeAmount,
			walletOrder.Currency,
			walletOrder.ExchangeAddress,
			walletOrder.DoneKey,
		)

		if settings != nil {
			for _, id := range settings.GetTelegramAdminsIds() {
				for _, botToken := range settings.GetTelegramNotificationChannelsTokens() {
					NotifyAdmins(text, id, botToken)
				}
			}
		}
	} else {
		return nil
	}

	return nil
}

func NotifyAdmins(text string, chatId int, botToken string) {
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	reqBody := &runtime.TelegramOutgoingMessage{
		ChatID:    uint(chatId),
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
