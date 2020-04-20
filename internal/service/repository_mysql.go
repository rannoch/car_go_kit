package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"strings"
)

type repositoryMysql struct {
	db *sqlx.DB
}

func NewRepositoryMysql(db *sqlx.DB) *repositoryMysql {
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	return &repositoryMysql{db: db}
}

func (r repositoryMysql) Store(ctx context.Context, car *Car) error {
	var query string

	if car.Id == 0 {
		query = `INSERT INTO car (brand, model, created) VALUES (?, ?, now())`

		result, err := r.db.ExecContext(ctx, query, car.Brand, car.Model)
		if err != nil {
			return err
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		car.Id = int(lastInsertId)
	} else {
		query = "UPDATE car SET brand = ?, model = ? WHERE id = ?"

		_, err := r.db.ExecContext(ctx, query, car.Brand, car.Model, car.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r repositoryMysql) Find(ctx context.Context, id int) (*Car, error) {
	var car Car

	err := r.db.GetContext(ctx, &car, "SELECT * FROM car WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	return &car, nil
}

func (r repositoryMysql) FindAll(ctx context.Context) ([]*Car, error) {
	var cars []*Car

	err := r.db.SelectContext(ctx, &cars, "SELECT * FROM car")

	return cars, err
}

func (r repositoryMysql) Del(ctx context.Context, id int) error {
	query := "DELETE FROM car WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
