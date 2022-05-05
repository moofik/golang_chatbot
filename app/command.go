package app

import (
	"bot-daedalus/bot/runtime"
	"bot-daedalus/models"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
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

	if cmd == "validate_btc_address" {
		validity := false

		if len(arguments) > 0 {
			validity = arguments[0].(bool)
		}
		// в метаданных cmd = text_input, потому что любая коммнада валидации ввода - подмножество комманды UserInput
		// это нужно для того чтобы можно было найти комманду по прототипу - прототип в данном случае комманда текстового ввода
		return &ValidateBtcAddress{
			Validity: validity,
			Text:     "",
			Metadata: &runtime.CommandMetadata{Cmd: "text_input", Place: place, Uniqueness: strconv.FormatBool(validity)},
		}
	}

	if cmd == "preorder_processing" {
		validity := false

		if len(arguments) > 0 {
			validity = arguments[0].(bool)
		}
		// в метаданных cmd = text_input, потому что любая коммнада валидации ввода - подмножество комманды UserInput
		// это нужно для того чтобы можно было найти комманду по прототипу - прототип в данном случае комманда текстового ввода
		return &PreorderProcessing{
			OrderRepository: &models.OrderRepository{DB: cr.DB},
			TokenRepository: &models.TokenRepository{DB: cr.DB},
			Validity:        validity,
			Text:            "",
			Metadata:        &runtime.CommandMetadata{Cmd: "recognize_input", Place: place, Uniqueness: strconv.FormatBool(validity) + "preorder_processing"},
		}
	}

	if cmd == "payment_processing" {
		validity := false

		if len(arguments) > 0 {
			validity = arguments[0].(bool)
		}
		// в метаданных cmd = text_input, потому что любая коммнада валидации ввода - подмножество комманды UserInput
		// это нужно для того чтобы можно было найти комманду по прототипу - прототип в данном случае комманда текстового ввода
		return &PaymentProcessing{
			OrderRepository: &models.OrderRepository{DB: cr.DB},
			TokenRepository: &models.TokenRepository{DB: cr.DB},
			Validity:        validity,
			Text:            "",
			Metadata:        &runtime.CommandMetadata{Cmd: "recognize_input", Place: place, Uniqueness: strconv.FormatBool(validity) + "payment_processing"},
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
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
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
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
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

type ValidateBtcAddress struct {
	Text     string
	Metadata *runtime.CommandMetadata
	Validity bool
}

func (c *ValidateBtcAddress) ToUniquenessHash() string {
	return runtime.ToUniquenessHash(c.Metadata)
}

func (c *ValidateBtcAddress) ToHash() string {
	return runtime.ToHash(c.Metadata)
}

func (c *ValidateBtcAddress) ToProtoHash() string {
	return runtime.ToProtoHash(c.Metadata)
}

func (c *ValidateBtcAddress) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *ValidateBtcAddress) GetMetadata() *runtime.CommandMetadata {
	return c.Metadata
}

func (c *ValidateBtcAddress) GetInput() string {
	return c.Text
}

func (c *ValidateBtcAddress) GetCaption() string {
	return "validate_btc_address"
}

func (c *ValidateBtcAddress) Pass(p runtime.ChatProvider, initCmd runtime.Command, t runtime.TokenProxy) (bool, error) {
	actualValidity := true

	if len(initCmd.GetInput()) > 35 {
		fmt.Printf("Validitation error detected: %s\n", "too long BTC address")
		actualValidity = false
		return c.Validity == actualValidity, nil
	}

	return c.Validity == actualValidity, nil
}

func (c *ValidateBtcAddress) GetType() string {
	return "validate_btc_address"
}

type PreorderProcessing struct {
	OrderRepository *models.OrderRepository
	TokenRepository *models.TokenRepository
	Text            string
	Metadata        *runtime.CommandMetadata
	Validity        bool
}

func (c *PreorderProcessing) ToUniquenessHash() string {
	return runtime.ToUniquenessHash(c.Metadata)
}

func (c *PreorderProcessing) ToHash() string {
	return runtime.ToHash(c.Metadata)
}

func (c *PreorderProcessing) ToProtoHash() string {
	return runtime.ToProtoHash(c.Metadata)
}

func (c *PreorderProcessing) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *PreorderProcessing) GetMetadata() *runtime.CommandMetadata {
	return c.Metadata
}

func (c *PreorderProcessing) GetInput() string {
	return c.Text
}

func (c *PreorderProcessing) GetCaption() string {
	return "preorder_processing"
}

func (c *PreorderProcessing) Pass(p runtime.ChatProvider, initCmd runtime.Command, t runtime.TokenProxy) (bool, error) {
	fmt.Printf("PreorderProcessing logging \n")
	input := initCmd.GetInput()

	if len(input) <= 8 {
		return false, nil
	}

	cmd := input[:8]
	doneKey := input[8:]
	order := c.OrderRepository.FindByDoneKey(doneKey)

	if order == nil {
		fmt.Printf("order not found with done key %s\n", doneKey)
		return false, nil
	}

	token := c.TokenRepository.FindById(order.TokenID)

	if token == nil {
		fmt.Printf("PaymentProcessing: token not found with id %s\n", order.TokenID)
		return false, nil
	}

	actualValidity := true

	if cmd == "/accept_" && order != nil {
		fmt.Printf("PreorderProcessing validity setting to true\n")
		actualValidity = true
	} else if cmd == "/refuse_" && order != nil {
		fmt.Printf("PreorderProcessing validity setting to false\n")
		actualValidity = false
	} else {
		fmt.Printf("PreorderProcessing: can't recognize command %s\n", cmd)
		return false, fmt.Errorf("PreorderProcessing: can't recognize command")
	}

	if actualValidity == false && actualValidity == c.Validity {
		pendingCmd := runtime.CreatePendingCommand("", "fail")
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
	}

	fmt.Printf("PreorderProcessing validity %b\n", c.Validity == actualValidity)
	return c.Validity == actualValidity, nil
}

func (c *PreorderProcessing) GetType() string {
	return "preorder_processing"
}

type PaymentProcessing struct {
	OrderRepository *models.OrderRepository
	TokenRepository *models.TokenRepository
	Text            string
	Metadata        *runtime.CommandMetadata
	Validity        bool
}

func (c *PaymentProcessing) ToUniquenessHash() string {
	return runtime.ToUniquenessHash(c.Metadata)
}

func (c *PaymentProcessing) ToHash() string {
	return runtime.ToHash(c.Metadata)
}

func (c *PaymentProcessing) ToProtoHash() string {
	return runtime.ToProtoHash(c.Metadata)
}

func (c *PaymentProcessing) Debug() string {
	return fmt.Sprintf("cmd: %s, state name: %s, hash: %s, CMD + PLACE + UNIQ: %s", c.Metadata.Cmd, c.Metadata.Place, c.ToHash(), c.Metadata.Cmd+c.Metadata.Place+c.Metadata.Uniqueness)
}

func (c *PaymentProcessing) GetMetadata() *runtime.CommandMetadata {
	return c.Metadata
}

func (c *PaymentProcessing) GetInput() string {
	return c.Text
}

func (c *PaymentProcessing) GetCaption() string {
	return "payment_processing"
}

func (c *PaymentProcessing) Pass(p runtime.ChatProvider, initCmd runtime.Command, t runtime.TokenProxy) (bool, error) {
	input := initCmd.GetInput()
	hasPayed, doneKeySuccess := c.ExtractOrderKey(input, "/payed_yes_")
	fmt.Printf("DONE KEY SUCCESS: %s\n", doneKeySuccess)
	hasNoPayed, doneKeyFail := c.ExtractOrderKey(input, "/payed_no_")
	fmt.Printf("DONE KEY FAIL: %s\n", doneKeyFail)
	var order *models.Order

	if hasPayed {
		order = c.OrderRepository.FindByDoneKey(doneKeySuccess)
	} else if hasNoPayed {
		order = c.OrderRepository.FindByDoneKey(doneKeyFail)
	}

	if order == nil {
		fmt.Printf("PaymentProcessing: order not found for input %s\n", input)
		return false, nil
	}

	token := c.TokenRepository.FindById(order.TokenID)

	if token == nil {
		fmt.Printf("PaymentProcessing: token not found with id %s\n", order.TokenID)
		return false, nil
	}

	actualValidity := true

	if hasPayed && order != nil {
		actualValidity = true
	} else if hasNoPayed && order != nil {
		actualValidity = false
	} else {
		fmt.Printf("can't recognize command %s\n", input)
		return false, nil
	}

	if actualValidity == true && actualValidity == c.Validity {
		order.IsDone = true
		pendingCmd := runtime.CreatePendingCommand("", "success")
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
		c.OrderRepository.Persist(order)
	}

	if actualValidity == false && actualValidity == c.Validity {
		order.IsDone = true
		pendingCmd := runtime.CreatePendingCommand("", "fail")
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
		c.OrderRepository.Persist(order)
	}

	return c.Validity == actualValidity, nil
}

func (c *PaymentProcessing) GetType() string {
	return "payment_processing"
}

func (c *PaymentProcessing) ExtractOrderKey(input, searchCmd string) (bool, string) {
	cmdIdx := strings.Index(input, searchCmd)
	if cmdIdx != 0 {
		fmt.Printf("PaymentProcessing: can't recognize command %s while extracting\n", input)
		return false, ""
	}
	doneIdx := cmdIdx + len(searchCmd)
	return true, input[doneIdx:]
}
