package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"test-server/models"
)

func FetchBTCRate(models.CoinRateApiOutput) (models.CoinRateApiOutput, error) {
	// Call to btc price API
	response, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice/USD.json")
	// check for erros
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		var empty models.CoinRateApiOutput
		return empty, err
	}

	// read the rsponse body
	data, _ := ioutil.ReadAll(response.Body)

	// read data and unmarashell data into CoinRateApiOutput struct
	var coinRate = new(models.CoinRateApiOutput)
	err = json.Unmarshal(data, &coinRate)
	// check for errors
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		var empty models.CoinRateApiOutput
		return empty, err
	}
	return *coinRate, nil
}
