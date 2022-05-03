package app

import (
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

//app package actions

type CommandRegistry struct {
	DB *gorm.DB
}

func (cr *CommandRegistry) CommandRegistryHandler(cmd string, place string, arguments []interface{}) runtime.Command {
	if cmd == "validate_wallet_order" {
		validity := false

		if len(arguments) > 0 {
			validity = arguments[0].(bool)
		}
		// в метаданных cmd = text_input, потому что любая коммнада валидации ввода - подмножество комманды UserInput
		// это нужно для того чтобы можно было найти комманду по прототипу - прототип в данном случае комманда текстового ввода
		return &ValidateWalletOrderCommand{
			Validity:         validity,
			WalletRepository: &models.WalletRepository{DB: cr.DB},
			Text:             "",
			Metadata:         &runtime.CommandMetadata{Cmd: "text_input", Place: place, Uniqueness: strconv.FormatBool(validity)},
		}
	}

	if cmd == "validate_market_order" {
		validity := false

		if len(arguments) > 0 {
			validity = arguments[0].(bool)
		}
		// в метаданных cmd = text_input, потому что любая коммнада валидации ввода - подмножество комманды UserInput
		// это нужно для того чтобы можно было найти комманду по прототипу - прототип в данном случае комманда текстового ввода
		return &ValidateMarketOrderCommand{
			Validity: validity,
			Text:     "",
			Metadata: &runtime.CommandMetadata{Cmd: "text_input", Place: place, Uniqueness: strconv.FormatBool(validity)},
		}
	}

	if cmd == "recognize_order" {
		validity := false

		if len(arguments) > 0 {
			validity = arguments[0].(bool)
		}
		// в метаданных cmd = text_input, потому что любая коммнада валидации ввода - подмножество комманды UserInput
		// это нужно для того чтобы можно было найти комманду по прототипу - прототип в данном случае комманда текстового ввода
		return &OrderConfirmation{
			OrderRepository: &models.OrderRepository{DB: cr.DB},
			TokenRepository: &models.TokenRepository{DB: cr.DB},
			Validity:        validity,
			Text:            "",
			Metadata:        &runtime.CommandMetadata{Cmd: "recognize_input", Place: place, Uniqueness: strconv.FormatBool(validity)},
		}
	}

	return nil
}

type ValidateWalletOrderCommand struct {
	*models.WalletRepository
	Text     string
	Metadata *runtime.CommandMetadata
	Validity bool
}

func (c *ValidateWalletOrderCommand) ToUniquenessHash() string {
	return runtime.ToUniquenessHash(c.Metadata)
}

func (c *ValidateWalletOrderCommand) ToHash() string {
	return runtime.ToHash(c.Metadata)
}

func (c *ValidateWalletOrderCommand) ToProtoHash() string {
	return runtime.ToProtoHash(c.Metadata)
}

func (c *ValidateWalletOrderCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, text: %s, hash: %s, data: %s", c.Metadata.Cmd, c.Metadata.Place, c.Text, c.ToHash(), c.GetInput())
}

func (c *ValidateWalletOrderCommand) GetMetadata() *runtime.CommandMetadata {
	return c.Metadata
}

func (c *ValidateWalletOrderCommand) GetInput() string {
	return c.Text
}

func (c *ValidateWalletOrderCommand) GetCaption() string {
	return "validate_wallet_order"
}

func (c *ValidateWalletOrderCommand) Pass(p runtime.ChatProvider, initCmd runtime.Command, t runtime.TokenProxy) (bool, error) {
	wallet := c.WalletRepository.FindWalletByTokenId(t.GetId())
	extras := t.GetExtras()
	var actualValidity bool
	requestedAmount, err := strconv.ParseFloat(initCmd.GetInput(), 64)

	if err != nil {
		return false, err
	}

	if extras["wallet_order_currency"] == "BTC" {
		actualValidity = wallet.BalanceBTC >= requestedAmount
	} else if extras["wallet_order_currency"] == "ETH" {
		actualValidity = wallet.BalanceETH >= requestedAmount
	} else if extras["wallet_order_currency"] == "BNB" {
		actualValidity = wallet.BalanceBNB >= requestedAmount
		fmt.Println(actualValidity)
	} else if extras["wallet_order_currency"] == "USDT" {
		actualValidity = wallet.BalanceUSDT >= requestedAmount
	} else {
		return false, nil
	}

	return c.Validity == actualValidity, nil
}

func (c *ValidateWalletOrderCommand) GetType() string {
	return "validate_wallet_order"
}

type ValidateMarketOrderCommand struct {
	Text     string
	Metadata *runtime.CommandMetadata
	Validity bool
}

func (c *ValidateMarketOrderCommand) ToUniquenessHash() string {
	return runtime.ToUniquenessHash(c.Metadata)
}

func (c *ValidateMarketOrderCommand) ToHash() string {
	return runtime.ToHash(c.Metadata)
}

func (c *ValidateMarketOrderCommand) ToProtoHash() string {
	return runtime.ToProtoHash(c.Metadata)
}

func (c *ValidateMarketOrderCommand) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, text: %s, hash: %s, data: %s", c.Metadata.Cmd, c.Metadata.Place, c.Text, c.ToHash(), c.GetInput())
}

func (c *ValidateMarketOrderCommand) GetMetadata() *runtime.CommandMetadata {
	return c.Metadata
}

func (c *ValidateMarketOrderCommand) GetInput() string {
	return c.Text
}

func (c *ValidateMarketOrderCommand) GetCaption() string {
	return "validate_wallet_order"
}

func (c *ValidateMarketOrderCommand) Pass(p runtime.ChatProvider, initCmd runtime.Command, t runtime.TokenProxy) (bool, error) {
	extras := t.GetExtras()
	actualValidity := true
	amt, err := strconv.ParseFloat(initCmd.GetInput(), 64)

	if err != nil {
		fmt.Println("Validitation error detected: %s", err.Error())
		actualValidity = false
		return c.Validity == actualValidity, nil
	}

	ps, _, _ := ConvertCrypto(extras["market_order_currency"], "RUB", amt, true)

	if ps < 1000 || ps > 15000 {
		actualValidity = false
	}

	return c.Validity == actualValidity, nil
}

func (c *ValidateMarketOrderCommand) GetType() string {
	return "validate_market_order"
}

type OrderConfirmation struct {
	OrderRepository *models.OrderRepository
	TokenRepository *models.TokenRepository
	Text            string
	Metadata        *runtime.CommandMetadata
	Validity        bool
}

func (c *OrderConfirmation) ToUniquenessHash() string {
	return runtime.ToUniquenessHash(c.Metadata)
}

func (c *OrderConfirmation) ToHash() string {
	return runtime.ToHash(c.Metadata)
}

func (c *OrderConfirmation) ToProtoHash() string {
	return runtime.ToProtoHash(c.Metadata)
}

func (c *OrderConfirmation) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, text: %s, hash: %s, data: %s", c.Metadata.Cmd, c.Metadata.Place, c.Text, c.ToHash(), c.GetInput())
}

func (c *OrderConfirmation) GetMetadata() *runtime.CommandMetadata {
	return c.Metadata
}

func (c *OrderConfirmation) GetInput() string {
	return c.Text
}

func (c *OrderConfirmation) GetCaption() string {
	return "recognize_input"
}

func (c *OrderConfirmation) Pass(p runtime.ChatProvider, initCmd runtime.Command, t runtime.TokenProxy) (bool, error) {
	fmt.Println("ORDER CONFIRMATION LOGGING 1")
	input := initCmd.GetInput()

	if len(input) <= 8 {
		return false, nil
	}

	fmt.Println(input)

	cmd := input[:8]
	doneKey := input[8:]
	order := c.OrderRepository.FindByDoneKey(doneKey)

	if order == nil {
		panic("PIZDEC")
	}

	fmt.Println("ORDER CONFIRMATION LOGGING 2")
	actualValidity := true
	var pendingCmd runtime.Command

	if cmd == "/accept_" && order != nil {
		actualValidity = true
		fmt.Println("ORDER CONFIRMATION LOGGING 3")
		pendingCmd = runtime.CreatePendingCommand("", "success")
	} else if cmd == "/refuse_" && order != nil {
		fmt.Println("ORDER CONFIRMATION LOGGING 4")
		pendingCmd = runtime.CreatePendingCommand("", "fail")

		actualValidity = false
	} else {
		fmt.Println("ORDER CONFIRMATION LOGGING 5")
		if order == nil {
			fmt.Printf("order not found with done key %s\n", doneKey)
			return false, fmt.Errorf("order not found with done key %s", doneKey)
		}

		fmt.Printf("can't recognize command %s\n", cmd)
		return false, fmt.Errorf("can't recognize command")
	}

	order.IsDone = true
	c.OrderRepository.Persist(order)
	token := c.TokenRepository.FindById(order.TokenID)
	actionRegistry := ActionRegistry{DB: c.TokenRepository.DB}
	commandRegistry := CommandRegistry{DB: c.TokenRepository.DB}
	bot := runtime.DefaultBot{
		ScenarioPath:       "config/scenario",
		ScenarioName:       "cryptobot",
		TokenFactory:       models.TokenFactory{DB: c.TokenRepository.DB},
		TokenRepository:    &models.TokenRepository{DB: c.TokenRepository.DB},
		SettingsRepository: &models.SettingsRepository{DB: c.TokenRepository.DB},
		ActionRegistry:     actionRegistry.ActionRegistryHandler,
		CommandRegistry:    commandRegistry.CommandRegistryHandler,
		StateErrorHandler:  CryptobotStateErrorHandler,
	}
	_, _, scenario := bot.GetBaseActors(&runtime.DefaultSerializedMessageFactory{Ctx: nil})
	currentState := scenario.GetCurrentState(token)
	innerToken := scenario.HandleCommand(pendingCmd, currentState, token)
	c.TokenRepository.Persist(innerToken)

	fmt.Println("JEEEEEG!) ALL GOOD :))))")

	return c.Validity == actualValidity, nil
}

func (c *OrderConfirmation) GetType() string {
	return "recognize_input"
}
