package impl

import (
	"context"
	"database/sql"
)

func (r *repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.Db.BeginTx(ctx, nil)
}

func (r *repository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *repository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}
