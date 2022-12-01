// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: transactions.sql

package indexerdb

import (
	"context"
	"database/sql"

	"go.vocdoni.io/dvote/types"
)

const countTxReferences = `-- name: CountTxReferences :one
SELECT COUNT(*) FROM tx_references
`

func (q *Queries) CountTxReferences(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTxReferences)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createTxReference = `-- name: CreateTxReference :execresult
INSERT INTO tx_references (
	hash, block_height, tx_block_index
) VALUES (
	?, ?, ?
)
`

type CreateTxReferenceParams struct {
	Hash         types.Hash
	BlockHeight  int64
	TxBlockIndex int64
}

func (q *Queries) CreateTxReference(ctx context.Context, arg CreateTxReferenceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createTxReference, arg.Hash, arg.BlockHeight, arg.TxBlockIndex)
}

const getLastTxReferences = `-- name: GetLastTxReferences :many
SELECT id, hash, block_height, tx_block_index FROM tx_references
ORDER BY id DESC
LIMIT ?
`

func (q *Queries) GetLastTxReferences(ctx context.Context, limit int32) ([]TxReference, error) {
	rows, err := q.db.QueryContext(ctx, getLastTxReferences, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TxReference
	for rows.Next() {
		var i TxReference
		if err := rows.Scan(
			&i.ID,
			&i.Hash,
			&i.BlockHeight,
			&i.TxBlockIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTxReference = `-- name: GetTxReference :one
SELECT id, hash, block_height, tx_block_index FROM tx_references
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetTxReference(ctx context.Context, id int64) (TxReference, error) {
	row := q.db.QueryRowContext(ctx, getTxReference, id)
	var i TxReference
	err := row.Scan(
		&i.ID,
		&i.Hash,
		&i.BlockHeight,
		&i.TxBlockIndex,
	)
	return i, err
}

const getTxReferenceByHash = `-- name: GetTxReferenceByHash :one
SELECT id, hash, block_height, tx_block_index FROM tx_references
WHERE hash = ?
LIMIT 1
`

func (q *Queries) GetTxReferenceByHash(ctx context.Context, hash types.Hash) (TxReference, error) {
	row := q.db.QueryRowContext(ctx, getTxReferenceByHash, hash)
	var i TxReference
	err := row.Scan(
		&i.ID,
		&i.Hash,
		&i.BlockHeight,
		&i.TxBlockIndex,
	)
	return i, err
}
