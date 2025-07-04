// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const allOrders = `-- name: AllOrders :many
SELECT
    number, status, accrual, uploaded_at, user_id
FROM
    orders
WHERE
    user_id = $1
ORDER BY
    uploaded_at DESC
`

func (q *Queries) AllOrders(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	rows, err := q.db.Query(ctx, allOrders, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.Number,
			&i.Status,
			&i.Accrual,
			&i.UploadedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createOrder = `-- name: CreateOrder :exec
INSERT INTO orders (
    number,
    user_id,
    status,
    accrual,
    uploaded_at
) VALUES (
    $1,
    $2,
    $3,
    0,
    NOW()
)
`

type CreateOrderParams struct {
	Number string
	UserID uuid.UUID
	Status string
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) error {
	_, err := q.db.Exec(ctx, createOrder, arg.Number, arg.UserID, arg.Status)
	return err
}

const getOrderForUpdate = `-- name: GetOrderForUpdate :one
SELECT
    number, status, accrual, uploaded_at, user_id
FROM
    orders
WHERE
    number = $1
FOR UPDATE
`

func (q *Queries) GetOrderForUpdate(ctx context.Context, number string) (Order, error) {
	row := q.db.QueryRow(ctx, getOrderForUpdate, number)
	var i Order
	err := row.Scan(
		&i.Number,
		&i.Status,
		&i.Accrual,
		&i.UploadedAt,
		&i.UserID,
	)
	return i, err
}

const increaseBalance = `-- name: IncreaseBalance :exec
UPDATE balances
SET
    balance = balance + COALESCE(orders.accrual, 0)
FROM
    orders
WHERE
    balances.user_id = orders.user_id
  AND
    orders.number = $1
`

func (q *Queries) IncreaseBalance(ctx context.Context, number string) error {
	_, err := q.db.Exec(ctx, increaseBalance, number)
	return err
}

const updateOrder = `-- name: UpdateOrder :exec
UPDATE orders
SET
    status = $1,
    accrual = $2
WHERE
    number = $3
`

type UpdateOrderParams struct {
	Status  string
	Accrual pgtype.Numeric
	Number  string
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) error {
	_, err := q.db.Exec(ctx, updateOrder, arg.Status, arg.Accrual, arg.Number)
	return err
}
