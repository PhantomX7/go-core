package database_transaction

import "gorm.io/gorm"

type Client interface {
	NewTransaction() *gorm.DB
}

type TransactionManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) Client {
	return &TransactionManager{
		db: db,
	}
}

func (tm *TransactionManager) NewTransaction() *gorm.DB {
	return tm.db.Begin()
}
