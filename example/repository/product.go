package repository

import (
	"context"
	"database/sql"

	"github.com/tungnt/goit/example/model"
)

type ProductRepo interface {
	GetOne(ctx context.Context, id int64) (*model.Product, error)
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) GetOne(ctx context.Context, id int64) (*model.Product, error) {
	query := "SELECT name, category_id, warehouse_id FROM product WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	var (
		name        string
		categoryID  int64
		warehouseID int64
	)
	err := row.Scan(&name, &categoryID, &warehouseID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if nil != err {
		return nil, err
	}

	return &model.Product{
		ID:          id,
		Name:        name,
		CategoryID:  categoryID,
		WarehouseID: warehouseID,
	}, nil
}
