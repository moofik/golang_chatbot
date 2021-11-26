package app

import "bot-daedalus/bot/runtime"

func CryptobotStateErrorHandler(p runtime.ChatProvider, ctx runtime.ProviderContext) {
	if ctx.Token.GetState() != "unknown" && ctx.Token.GetState() != "start" {
		_ = p.SendMarkupMessage(
			[]string{"Маркет💵", "Кошелек💠"},
			"К сожалению я не знаю такой комманды. Вы можете воспользоваться меню ниже.",
			ctx,
		)
	} else {
		_ = p.SendTextMessage(
			"Для взаимодействия с ботом вам необходимо перейти в меню.",
			ctx,
		)
	}
}
