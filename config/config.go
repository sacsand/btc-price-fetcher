package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// dir path for build
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// InitConfig
func InitConfig() {
	envLoader()
	// jobs.CoinPriceFeeder()
}

func envLoader() {
	// load .env file
	notLocalDevelop := os.Getenv("ENV") != ""
	if notLocalDevelop {
		return
	}
	fmt.Println(filepath.Dir(basepath), "filepath.Dir(basepath)")
	err := godotenv.Load(filepath.Dir(basepath) + "/.env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
}
