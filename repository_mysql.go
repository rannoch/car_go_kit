package car

import (
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

func (r repositoryMysql) Store(car *Car) error {
	var query string

	if car.Id == 0 {
		query = `INSERT INTO car (brand, model, created) VALUES (?, ?, now())`

		result, err := r.db.Exec(query, car.Brand, car.Model)
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

		_, err := r.db.Exec(query, car.Brand, car.Model, car.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r repositoryMysql) Find(id int) (*Car, error) {
	var car Car

	err := r.db.Get(&car, "SELECT * FROM car WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	return &car, nil
}

func (r repositoryMysql) FindAll() ([]*Car, error) {
	var cars []*Car

	err := r.db.Select(&cars, "SELECT * FROM car")

	return cars, err
}

func (r repositoryMysql) Del(id int) error {
	query := "DELETE FROM car WHERE id = ?"

	_, err := r.db.Exec(query, id)
	return err
}

