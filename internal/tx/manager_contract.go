package tx

import "context"

type TransactionManager interface {
	Begin(parent context.Context, withTx bool, fn func(ctx context.Context) error) error
	GetConnection(ctx context.Context) (DB, error)
}
