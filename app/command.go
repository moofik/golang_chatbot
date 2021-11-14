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
	fmt.Println("JEEEEEEEEH)))")
	wallet := c.WalletRepository.FindWalletByTokenId(t.GetId())
	extras := t.GetExtras()
	var actualValidity bool
	requestedAmount, err := strconv.ParseFloat(initCmd.GetInput(), 64)

	if err != nil {
		fmt.Println("error converting amt")
		return false, err
	}

	fmt.Println("ct amt : " + string(extras["wallet_order_currency"]))
	s := fmt.Sprintf("%f", requestedAmount)
	fmt.Println("rq amt : " + s)
	fmt.Println("wallet" + fmt.Sprintf("%d", wallet.ID))

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
