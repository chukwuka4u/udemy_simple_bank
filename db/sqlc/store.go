package db

import (
	"fmt"
	"database/sql"
	"context"
)
// store provies all functions to execute queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db: db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %s, rebErr: %s", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json: from_account_id`
	ToAccountID int64 `json: to_ccount_id`
	Amount int64 `json: amount`
}

type TransferTxResult struct {
	Transfer Transfer `json: transfer`
	FromAccount Account `json: from_account`
	ToAccount Account `json: to_account`
	ToEntry Entry `json: to_entry`
	FromEntry Entry `json: from_entry`
}

var txKey = struct{}{}

//creates a transfer record, add account entries, update account balance, within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		txName := ctx.Value(txKey)
		var err error
		
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		
		
		fmt.Println(txName, "update account 1")
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "update account 2")
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}