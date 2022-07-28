package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ajikamaludin/api-raya-ojt/app/configs"
	"github.com/ajikamaludin/api-raya-ojt/app/models"
	"github.com/ajikamaludin/api-raya-ojt/pkg/googlepubsub"
	"github.com/ajikamaludin/api-raya-ojt/pkg/gormdb"
	"github.com/ajikamaludin/api-raya-ojt/pkg/utils/constants"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(os.Getwd())
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// NOTE: maybe in future need more worker to execute transaction task thinks later
	subscriberName := constants.TRANSACTION_SUBSCRIBER_NAME
	configs := configs.GetInstance()
	googlepubsub := googlepubsub.New(configs.GooglePubSubConfig.ProjectName)
	fmt.Println("[SERVICES: STARTING]: Ok")
	googlepubsub.CreateTopicIfNotExists(constants.TRANSACTION_TOPIC_NAME)
	fmt.Println("[SERVICES: CREATE TOPIC]: Ok")
	googlepubsub.CreateSubscribtion(subscriberName)
	fmt.Println("[SERVICES: CREATE SUBSCRIBTION]: Ok")
	fmt.Println("[SERVICES: LISTENING]: Listen")
	err = googlepubsub.Subscibe(subscriberName, func(ctx context.Context, message *pubsub.Message) {
		// TODO: process transaction in here make call some function of anything
		message.Ack()
		id := string(message.Data)
		fmt.Println("Receive Transaction ID :", id)
		go ProcessTransaction(id)
	})
	if err != nil {
		fmt.Println(err)
		panic("Can't Create Subscribtion")
	}
}

func ProcessTransaction(id string) error {
	gormdb := gormdb.New()
	db, err := gormdb.GetInstance()
	defer gormdb.Conn.Close()
	fmt.Println("Process Transction ID : ", id)

	// Find Transaction In DB
	var transaction models.BankTransaction
	err = db.Preload("Bank").Find(&transaction, "id = ?", id).Error
	if err != nil {
		fmt.Println("Transaction ID:", id, " Error:", err)
		return err
	}

	// NOTE: call api bank , mock it with time.Sleep , and automated success
	time.Sleep(10 * time.Second)
	if transaction.Bank.Code == constants.RAYA_BANK_CODE {
		// Internal Transfer to Account
		// Create Transaction For Destination Account
		txdb := db.Begin()
		err := txdb.Model(&transaction).Update("status", constants.TRX_SUCCESS).Error
		if err != nil {
			fmt.Println("Transaction ID: ", id, " Error:", err)
			return err
		}

		var account models.Account
		err = txdb.Find(&account, "id = ?", transaction.AccountId).Error
		if err != nil {
			fmt.Println("Transaction ID: ", id, " Error:", err)
			return err
		}

		receivetrx := &models.BankTransaction{
			AccountId: account.ID,
			BankId:    transaction.BankId,
			UserId:    account.UserId,
			Debit:     transaction.Credit,
			Credit:    0,
			Status:    constants.TRX_SUCCESS,
		}

		err = txdb.Create(receivetrx).Error
		if err != nil {
			txdb.Rollback()
			fmt.Println("Transaction ID: ", id, " Error:", err)
			return err
		}

		err = txdb.Model(&account).Update("balance", account.Balance+receivetrx.Debit).Error
		if err != nil {
			txdb.Rollback()
			fmt.Println("Transaction ID: ", id, " Error:", err)
			return err
		}

		txdb.Commit()
	} else {
		// Extenal Transfer Account
		// Call Api Get Response , set success or fail
		db.Model(&transaction).Update("status", constants.TRX_SUCCESS)
	}

	fmt.Println("Done Process Transction ID : ", id)
	return nil
}
