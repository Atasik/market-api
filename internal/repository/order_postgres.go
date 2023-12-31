package repository

import (
	"fmt"
	"market/internal/model"
	"market/pkg/database/postgres"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderPostgresqlRepository struct {
	db *sqlx.DB
}

func NewOrderPostgresqlRepo(db *sqlx.DB) *OrderPostgresqlRepository {
	return &OrderPostgresqlRepository{db: db}
}

func (repo *OrderPostgresqlRepository) Create(cartID, userID int, order model.Order) (int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	defer func() (int, error) {
		if err := tx.Rollback(); err != nil {
			return 0, postgres.ParsePostgresError(err)
		}
		return 0, nil
	}()

	query := fmt.Sprintf("INSERT INTO %s (created_at, delivered_at, user_id) VALUES ($1, $2, $3) RETURNING id", ordersTable)
	row := tx.QueryRow(query, order.CreatedAt, order.DeliveredAt, userID)
	if err = row.Scan(&order.ID); err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	if len(order.Products) != 0 {
		var insertQueryBuilder strings.Builder

		insertQueryBuilder.WriteString(fmt.Sprintf("INSERT INTO %s (order_id, product_id, purchased_amount)", productsOrdersTable))

		args := []interface{}{}
		argID := 1
		for _, prod := range order.Products {
			args = append(args, order.ID, prod.ID, prod.PurchasedAmount)
			insertQueryBuilder.WriteString(fmt.Sprintf(`($%d,$%d,$%d),`, argID, argID+1, argID+2)) //nolint:gomnd
			argID += 3
		}

		query = strings.TrimSuffix(insertQueryBuilder.String(), ",")

		if _, err = tx.Exec(query, args...); err != nil {
			return 0, postgres.ParsePostgresError(err)
		}
	}

	query = fmt.Sprintf(`UPDATE %s AS p
						 SET amount = p.amount - pc.purchased_amount
					     FROM %s AS pc 
						 WHERE pc.product_id = p.id AND pc.cart_id = $1`, productsTable, productsCartsTable)
	if _, err = tx.Exec(query, cartID); err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE cart_id = $1`, productsCartsTable)
	if _, err = tx.Exec(query, cartID); err != nil {
		return 0, postgres.ParsePostgresError(err)
	}

	return order.ID, postgres.ParsePostgresError(tx.Commit())
}

func (repo *OrderPostgresqlRepository) GetAll(userID int, q model.OrderQueryInput) ([]model.Order, error) {
	var orders []model.Order
	query := fmt.Sprintf(`SELECT o.id, o.created_at, o.delivered_at FROM %s o
			              INNER JOIN %s u on o.user_id = u.id
			              WHERE u.id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3`, ordersTable, usersTable, q.SortBy, q.SortOrder)

	if err := repo.db.Select(&orders, query, userID, q.Limit, q.Offset); err != nil {
		return []model.Order{}, postgres.ParsePostgresError(err)
	}

	return orders, nil
}

func (repo *OrderPostgresqlRepository) GetByID(orderID int) (model.Order, error) {
	var order model.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", ordersTable)

	if err := repo.db.Get(&order, query, orderID); err != nil {
		return model.Order{}, postgres.ParsePostgresError(err)
	}

	return order, nil
}

func (repo *OrderPostgresqlRepository) GetProductsByOrderID(orderID int, q model.ProductQueryInput) ([]model.Product, error) {
	var products []model.Product
	query := fmt.Sprintf(`SELECT p.id, p.user_id, p.title, p.price, p.tag, p.category, p.description, p.amount, p.created_at, p.updated_at, p.views, p.image_url FROM %s p 
			              INNER JOIN %s po on po.product_id = p.id
			              INNER JOIN %s o on po.order_id = o.id
			              WHERE o.id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3`, productsTable, productsOrdersTable, ordersTable, q.SortBy, q.SortOrder)

	if err := repo.db.Select(&products, query, orderID, q.Limit, q.Offset); err != nil {
		return []model.Product{}, postgres.ParsePostgresError(err)
	}

	return products, nil
}
