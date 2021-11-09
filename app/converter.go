package app

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

const COIN_API_KEY = "14F07201-05F6-40DE-8A4D-00DBE5BCAE6A"
const COIN_MARKET_API_KEY = "4de8f1fc-2780-4e8f-921b-9ce0a247ae36"
const MONEY_BUY_COEFFICIENT = 1.03  // 3%
const MONEY_SELL_COEFFICIENT = 0.97 // 3%
const CURRENCY_BTC = "BTC"
const CURRENCY_ETH = "ETH"
const CURRENCY_USDT = "USDT"
const CURRENCY_BNB = "BNB"

type CoinApiRateResponse struct {
	Time         string
	AssetIdBase  string
	AssetIdQuote string
	Rate         float64
}

type CoinMarketRateResponse struct {
	Data []struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  struct {
			Rub struct {
				Price float64 `json:"price"`
			} `json:"RUB"`
		} `json:"quote"`
	} `json:"data"`
}

type CoinGeckoRatesResponse struct {
	Ethereum struct {
		Rub float64 `json:"rub"`
	} `json:"ethereum"`
	Bitcoin struct {
		Rub float64 `json:"rub"`
	} `json:"bitcoin"`
	Tether struct {
		Rub float64 `json:"rub"`
	} `json:"tether"`
	BinanceCoin struct {
		Rub float64 `json:"rub"`
	} `json:"binancecoin"`
}

func ConvertCrypto(from string, to string, amount float64, buy bool) (int, int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids=ethereum%2Cbitcoin%2Ctether%2Cbinancecoin&vs_currencies=rub", nil)
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("Accepts", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		return 0, 0, err
	}
	defer resp.Body.Close()

	apiResponse := CoinGeckoRatesResponse{}
	json.NewDecoder(resp.Body).Decode(&apiResponse)
	price := 0.0

	if from == CURRENCY_BTC {
		price = apiResponse.Bitcoin.Rub
	}
	if from == CURRENCY_ETH {
		price = apiResponse.Ethereum.Rub
	}
	if from == CURRENCY_USDT {
		price = apiResponse.Tether.Rub
	}
	if from == CURRENCY_BNB {
		price = apiResponse.BinanceCoin.Rub
	}

	res := 0.0
	actualPrice := price * amount

	if buy {
		res = price * MONEY_BUY_COEFFICIENT * amount
	} else {
		res = price * MONEY_SELL_COEFFICIENT * amount
	}

	return int(math.Ceil(res)), int(math.Ceil(actualPrice)), nil
}
