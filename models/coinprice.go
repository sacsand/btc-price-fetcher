package models

import "time"

// DB Module CoinPrice
type CoinPrice struct {
	ID   int `gorm:"primaryKey"`
	Rate float32
	Time time.Time
}

/**
OUT PUT OF BTCUSD Rate API
https://api.coindesk.com/v1/bpi/currentprice/USD.json

{
	"time": {
		"updated": "Jul 6, 2021 15:22:00 UTC",
		"updatedISO": "2021-07-06T15:22:00+00:00",
		"updateduk": "Jul 6, 2021 at 16:22 BST"
	},
	"disclaimer": "This data was produced from the CoinDesk Bitcoin Price Index (USD). Non-USD currency data converted using hourly conversion rate from openexchangerates.org",
	"bpi": {
		"USD": {
			"code": "USD",
			"rate": "33,997.6900",
			"description": "United States Dollar",
			"rate_float": 33997.69
		}
	}
	}
*/

// Struct to match above api response
type CoinRateApiOutput struct {
	Time TimeObj `json:"time"`
	Bpi  Bpi     `json:"bpi"`
}

type TimeObj struct {
	Updated    string `json:"updated"`
	UpdatedISO string `json:"updatedISO"`
}
type Bpi struct {
	Usd UsdObj `json:"USD"`
}

type UsdObj struct {
	Rate float32 `json:"rate_float"`
}
