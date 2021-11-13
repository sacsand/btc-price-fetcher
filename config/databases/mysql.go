package databases

import (
	"fmt"
	"os"
	"sync"
	"test-server/config"
	"test-server/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var once sync.Once

// DriverMySQL - MySQL connection grouping
type DriverMySQL struct {
	conn *gorm.DB
}

var db *DriverMySQL

// GetDB - Create Db singleton once, return same singleton db after, default connection,
func InitMySql() *gorm.DB {
	config.InitConfig()
	once.Do(func() {
		fmt.Println("Mysql singleton creation.")
		db = &DriverMySQL{conn: ConnectToMysql()}
	})
	return db.conn
}

// ConnectDB connect to db
func ConnectToMysql() *gorm.DB {

	var err error

	dsn := fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC`,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	if err != nil {
		panic(err)
	}

	dbconn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection open to database")

	dbconn.AutoMigrate(&models.CoinPrice{})

	db, _ := dbconn.DB()
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	fmt.Println("Database Migrated")

	return dbconn
}
