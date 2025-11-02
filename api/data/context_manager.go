package data

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

const (
	ctxTransactionKey string = "TRANSACTION"
)

func GetContextWithTransaction(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxTransactionKey, db)
}

func GetTransaction(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(ctxTransactionKey).(*gorm.DB); ok {
		return tx
	}
	log.Panicln(errors.New("error fetching transaction"))
	return nil
}
