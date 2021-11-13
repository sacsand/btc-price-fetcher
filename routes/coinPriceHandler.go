package routes

import (
	"fmt"
	"math/big"
	"test-server/pkg/coinprice"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Route Init
func CoinPriceRouter(api fiber.Router, coinPriceService coinprice.Service) {
	api.Get("/test-server/lastprice", getLastPrice(coinPriceService))
	api.Get("/test-server/:time", getRateByTime(coinPriceService))
	api.Get("/test-server/:from/:to", getAverageRate(coinPriceService))
}

// API For get last price
func getLastPrice(coinPriceService coinprice.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, _ := coinPriceService.FindLast()
		fmt.Println(data)
		return c.JSON(fiber.Map{
			"data":    data,
			"success": true,
		})
	}
}

// API to get price of given timestamp .
func getRateByTime(coinPriceService coinprice.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// final output struct
		type Output struct {
			Rate string
		}

		var output Output

		// 1) pass the UST time
		t1, e := time.Parse("2006-01-02T15:04:05Z07:00",
			c.Params("time"))
		if e != nil {
			fmt.Println(e)
		}

		t1.Format("2006-01-02 15:04:05.000") // re fromat time to match the databse record

		if e != nil {
			fmt.Println(e)
		}
		//  find the given time btc rate
		data, err := coinPriceService.FindByTime(t1.String())
		if err != nil {
			return c.JSON(fiber.Map{
				"data":    "wrong time format in time param,should be in 2006-01-02T15:04:05Z07:00 format",
				"success": false,
			})
		}

		//  if not found
		if data.ID == 0 {

			to := t1.Add(1 * time.Minute)

			from := t1.Add(-1 * time.Minute)
			fmt.Println(t1)

			avg, err := CalAverage(coinPriceService, from.String(), to.String())

			if !err {
				return c.JSON(fiber.Map{
					"data":    "Error calculating avg",
					"success": false,
				})
			}

			output.Rate = avg.String()
			fmt.Println("#######")
			fmt.Println(avg.String())
			fmt.Println("#######")

		} else {
			output.Rate = fmt.Sprintf("%v", data.Rate)
		}

		// return empty record
		fmt.Println(data)
		return c.JSON(fiber.Map{
			"data":    output,
			"success": true,
		})
	}
}

func getAverageRate(coinPriceService coinprice.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		type Output struct {
			Average string
		}

		// pass the UST time
		from, e := time.Parse("2006-01-02T15:04:05Z07:00",
			c.Params("from"))
		if e != nil {
			return c.JSON(fiber.Map{
				"data":    "wrong time format in time param,should be in 2006-01-02T15:04:05Z07:00 format",
				"success": false,
			})
		}

		to, e := time.Parse("2006-01-02T15:04:05Z07:00",
			c.Params("to"))
		if e != nil {
			return c.JSON(fiber.Map{
				"data":    "wrong time format in time param,should be in 2006-01-02T15:04:05Z07:00 format",
				"success": false,
			})
		}

		// check pair is valid
		validPair := from.Before(to)
		if !validPair {
			return c.JSON(fiber.Map{
				"data":    "not a valid time pair,should From before To",
				"success": false,
			})
		}
		// convert time format for  db save
		from.Format("2006-01-02 15:04:05.000") // re fromat time to match the databse record
		to.Format("2006-01-02 15:04:05.000")   // re fromat time to match the databse record

		// calculate average of the bitcoin for given time period
		avg, err := CalAverage(coinPriceService, from.String(), to.String())

		if !err {
			return c.JSON(fiber.Map{
				"data":    "Error calculating avg",
				"success": false,
			})
		}

		var output = Output{
			Average: avg.String(),
		}

		return c.JSON(fiber.Map{
			"data":    output,
			"success": true,
		})
	}
}

// Function to calcualte averagae
func CalAverage(coinPriceService coinprice.Service, from string, to string) (big.Float, bool) {
	var avg big.Float
	// fetch data
	data, _ := coinPriceService.FindBetween(from, to)

	sum := big.NewFloat(0)
	count := 0

	// find the sum of the rate on given range
	for _, v := range data {
		count++
		x := big.NewFloat(float64(v.Rate))
		sum.Add(sum, x)
	}
	// Error cehcking
	if count == 0 || sum == big.NewFloat(0) {
		return avg, false
	}

	kcount := big.NewFloat(float64(count))
	// cal avg
	avg.Quo(sum, kcount)

	return avg, true

}
