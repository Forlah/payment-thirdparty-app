package database

import "thirdparty-service/models"

//go:generate mockgen -source=mongodbstore.go -destination=../mocks/mongodbstore_mock.go -package=mocks
type MongoDBStore interface {
	GetAccountByID(accountId string) (*models.Account, error)
	UpdateAccountBalance(accountId string, amount float64) error
	CreateTransaction(transaction *models.Transaction) error
	GetPaymentByReferenceId(referenceId string) (*models.Transaction, error)
}
