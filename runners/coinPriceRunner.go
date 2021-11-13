package jobs

import (
	"fmt"
	"test-server/adapters"
	"test-server/models"
	"test-server/pkg/coinprice"
	"time"

	"github.com/robfig/cron/v3"
)

func CoinPriceRunner(coinPriceService coinprice.Service) {
	// creat new cron Obj
	cron := cron.New()
	// add new cron job
	cron.AddFunc("@every 0h0m60s", func() {
		fmt.Println("Btc price updated - (Every 1 Min)")

		// call to fetchBTC rate API
		response, err := adapters.FetchBTCRate(models.CoinRateApiOutput{})

		if err != nil {
			println("error fetching btc price")
		}

		// construct time
		t1, e := time.Parse("2006-01-02T15:04:05Z07:00",
			response.Time.UpdatedISO)
		if e != nil {
			fmt.Println(e)
		}

		// construc the coinPrice struct
		var currentPrice = models.CoinPrice{
			Rate: response.Bpi.Usd.Rate,
			Time: t1,
		}
		// save the rate to database
		coinPriceService.Save(currentPrice)

	})

	// start corn job
	cron.Start()
}
