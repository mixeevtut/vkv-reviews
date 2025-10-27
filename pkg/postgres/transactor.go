package postgres

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

// TODO добавить поддержку sqx.TxOptions

// WithinTransaction открываем транзакцию если ее нет
func (t *Transactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	var tx *sql.Tx
	var isOwnerTx bool
	committed := false

	if existTx, ok := GetTx(ctx); ok {
		isOwnerTx = false
		tx = existTx
	} else {
		newTx, err := t.db.BeginTx(ctx, nil)
		if err != nil {
			return errors.Wrap(err, "error beginning transaction")
		}
		isOwnerTx = true
		tx = newTx
	}

	defer func() {
		if isOwnerTx && !committed {
			if err := tx.Rollback(); err != nil {
				log.Printf("error rolling back transaction: %v", err)
			}
		}
	}()

	err := fn(injectTx(ctx, tx))
	if err != nil {
		return errors.Wrap(err, "error executing transaction")
	}

	if isOwnerTx {
		if err := tx.Commit(); err != nil {
			return errors.Wrap(err, "error committing transaction")
		}
		committed = true
	}

	return nil
}

type txKey struct{}

// injectTx внедряет транзакцию в контекст
func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx извлекает транзакцию из контекста
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}
