package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"test-server/config/databases"
	"test-server/models"
	"test-server/pkg/coinprice"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	Setup()                                  // setup fiber app
	mySqlConnection := databases.InitMySql() // create mysql connection

	coinPriceRepository := coinprice.NewRepo(mySqlConnection)     // Init coinPrice Repository
	coinPriceService := coinprice.NewService(coinPriceRepository) // Init coinPrice Service

	// delete the date in table for testing
	coinPriceService.Truncate()

	// feed data
	jsonFile, err := os.Open("_seed/coinprice.seed.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read josn content
	byteValue, _ := ioutil.ReadAll(jsonFile)

	seedData := []models.CoinPrice{}
	// unmarasahl data
	json.Unmarshal([]byte(byteValue), &seedData)
	// batch insert json data for testing
	coinPriceService.BatchInsert(seedData)

}

// Http Test
func TestAPI(t *testing.T) {

	// Setup the app as it is done in the main function
	app := Setup()

	// Define a structure for specifying input and output data of a single test case(Test Map)
	tests := []struct {
		description string

		// Test input
		route   string
		method  string
		payload *strings.Reader

		// Expected output
		partiallyAssert bool
		expectedError   bool
		expectedCode    int
		expectedBody    string
	}{
		{
			description:   "Test for non existing route",
			route:         "/not-exist",
			expectedError: false,
			expectedCode:  401,
			expectedBody:  "Not Found http://",
			method:        "GET",
			payload:       strings.NewReader(""),
		},
		{
			description:   "Test for get last price",
			route:         "/api/test-server/lastprice",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "{\"data\":{\"ID\":207,\"Rate\":33308.5,\"Time\":\"2021-07-08T03:32:00Z\"},\"success\":true}",
			method:        "GET",
			payload:       strings.NewReader(""),
		},
		{
			description:   "Test for get price for given timestamp",
			route:         "/api/test-server/2021-07-07T17:35:00+00:00",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "{\"data\":{\"Rate\":\"34517.96\"},\"success\":true}",
			method:        "GET",
			payload:       strings.NewReader(""),
		},
		{
			description:   "Test for get price for given timestamp but not found in DB, return Average price instead",
			route:         "/api/test-server/2021-07-07T17:35:20+00:00",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "{\"data\":{\"Rate\":\"34532.89453\"},\"success\":true}",
			method:        "GET",
			payload:       strings.NewReader(""),
		},
		{
			description:   "Test Find Average for given time range",
			route:         "/api/test-server/2021-07-07T17:34:00+00:00/2021-07-08T03:32:00+00:00",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "{\"data\":{\"Average\":\"33952.54272\"},\"success\":true}",
			method:        "GET",
			payload:       strings.NewReader(""),
		},
		{
			description:   "Test Find Average with wrong date format",
			route:         "/api/test-server/2021-07-07Twronfdate17:34:00+00:00/2021-07-08T03:32:00+00:00",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "{\"data\":\"wrong time format in time param,should be in 2006-01-02T15:04:05Z07:00 format\",\"success\":false}",
			method:        "GET",
			payload:       strings.NewReader(""),
		},
		{
			description:   "Test Find Average with wrong pair",
			route:         "/api/test-server/2021-07-08T03:32:00+00:00/2021-07-07T17:34:00+00:00",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "{\"data\":\"not a valid time pair,should From before To\",\"success\":false}",
			method:        "GET",
			payload:       strings.NewReader(""),
		},
	}

	// Iterate through test single test cases
	for _, test := range tests {

		// Create a new http request with the route
		// From the test case
		req, _ := http.NewRequest(
			test.method,
			test.route,
			test.payload,
		)

		req.Header.Add("Content-Type", "application/json")
		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		if test.partiallyAssert {
			// fmt.Println(string(body), "string(body)")
			assert.True(t, strings.Contains(string(body), test.expectedBody), test.description)
		} else {
			// Verify, that the reponse body equals the expected body
			assert.Equalf(t, test.expectedBody, string(body), test.description)
		}
	}
}
