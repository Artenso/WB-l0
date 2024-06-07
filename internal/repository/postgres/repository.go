package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Artenso/wb-l0/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	table        = "orders"
	orderUIDCol  = "order_uid"
	orderJSONCol = "order_json"
)

// IRepository working with repository
type IRepository interface {
	AddOrder(ctx context.Context, order *model.Order) error
	GetOrders(ctx context.Context) (pgx.Rows, error)
}

type repository struct {
	conn *pgxpool.Pool
}

// New reates new repository
func New(conn *pgxpool.Pool) IRepository {
	return &repository{
		conn: conn,
	}
}

// AddOrder adds order to db
func (r *repository) AddOrder(ctx context.Context, order *model.Order) error {
	jsonOrder, err := json.MarshalIndent(order, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal order: %s", err.Error())
	}

	builder := sq.Insert(table).
		Columns(orderUIDCol, orderJSONCol).
		Values(order.Order_uid, jsonOrder).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %s", err.Error())
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	rows.Close()

	return nil

}

// GetOrders gets all orders from db
func (r *repository) GetOrders(ctx context.Context) (pgx.Rows, error) {
	builder := sq.Select(orderJSONCol).
		From(table).
		PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %s", err.Error())
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
