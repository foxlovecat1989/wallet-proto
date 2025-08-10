package tx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"wallet-user-svc/pkg/utils/cx"
)

// TxWrapper wraps a database transaction and provides helper methods
type TxWrapper struct {
	tx *sqlx.Tx
}

// NewTxWrapper creates a new transaction wrapper
func NewTxWrapper(tx *sqlx.Tx) *TxWrapper {
	return &TxWrapper{tx: tx}
}

// GetTx returns the underlying transaction
func (tw *TxWrapper) GetTx() *sqlx.Tx {
	return tw.tx
}

// GetTxFromContext retrieves a transaction from context
func GetTxFromContext(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(cx.TransactionContextKey).(*sqlx.Tx)
	return tx, ok
}

// TransactionManager manages database transactions
type TransactionManager struct {
	db *sqlx.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(*TxWrapper) error) error {
	return tm.WithTransactionOptions(ctx, fn, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
}

// WithTransactionOptions executes a function within a database transaction with custom options
func (tm *TransactionManager) WithTransactionOptions(ctx context.Context, fn func(*TxWrapper) error, opts *sql.TxOptions) error {
	tx, err := tm.db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	txWrapper := NewTxWrapper(tx)

	// Execute the function
	if err := fn(txWrapper); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			// Log rollback error but return original error
			// You might want to add logging here: log.Printf("rollback failed: %v", rbErr)
			return err // return original error, not rollback error
		}
		return err
	}

	// Commit on success
	return tx.Commit()
}

// WithTransactionIsolation executes a function within a database transaction with specific isolation level
func (tm *TransactionManager) WithTransactionIsolation(ctx context.Context, fn func(*TxWrapper) error, isolation sql.IsolationLevel) error {
	return tm.WithTransactionOptions(ctx, fn, &sql.TxOptions{
		Isolation: isolation,
		ReadOnly:  false,
	})
}

// WithReadOnlyTransaction executes a function within a read-only database transaction
func (tm *TransactionManager) WithReadOnlyTransaction(ctx context.Context, fn func(*TxWrapper) error) error {
	return tm.WithTransactionOptions(ctx, fn, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
}

// WithSerializableTransaction executes a function within a serializable transaction
func (tm *TransactionManager) WithSerializableTransaction(ctx context.Context, fn func(*TxWrapper) error) error {
	return tm.WithTransactionIsolation(ctx, fn, sql.LevelSerializable)
}

// WithRepeatableReadTransaction executes a function within a repeatable read transaction
func (tm *TransactionManager) WithRepeatableReadTransaction(ctx context.Context, fn func(*TxWrapper) error) error {
	return tm.WithTransactionIsolation(ctx, fn, sql.LevelRepeatableRead)
}

// WithReadUncommittedTransaction executes a function within a read uncommitted transaction
func (tm *TransactionManager) WithReadUncommittedTransaction(ctx context.Context, fn func(*TxWrapper) error) error {
	return tm.WithTransactionIsolation(ctx, fn, sql.LevelReadUncommitted)
}
