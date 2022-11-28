package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CustomerDB struct {
	gorm.Model
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	AccountId *uint32   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type AccountDB struct {
	gorm.Model
	ID         uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Balance    uint64
	CustomerId uint32
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type TransactionsDB struct {
	gorm.Model
	ID                uint32 `gorm:"primary_key;auto_increment" json:"id"`
	DebitedAccountId  uint32
	CreditedAccountId uint32
	Value             uint64
}

var DBConnection *gorm.DB

func Connect() (*gorm.DB, error) {

	envVariables := []string{
		os.Getenv("POSTGRES_HOSTNAME"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"),
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", envVariables[0], envVariables[1], envVariables[2], envVariables[3], envVariables[4])
	DBConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
		return DBConnection, err
	}

	DBConnection.AutoMigrate(CustomerDB{}, AccountDB{}, TransactionsDB{})

	return DBConnection, nil
}

func Connection() (*gorm.DB, error) {
	if DBConnection != nil {
		return DBConnection, nil
	}

	envVariables := []string{
		os.Getenv("POSTGRES_HOSTNAME"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"),
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", envVariables[0], envVariables[1], envVariables[2], envVariables[3], envVariables[4])
	DBConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
		return DBConnection, err
	}

	return DBConnection, nil
}
