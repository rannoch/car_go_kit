package service

import (
	"context"
	"sync"
)

type repositoryInmem struct {
	cars          sync.Map
	autoIncrement int
}

func NewCarRepositoryInmem() *repositoryInmem {
	return &repositoryInmem{}
}

func (c *repositoryInmem) Store(_ context.Context, car *Car) error {
	if car.Id == 0 {
		c.autoIncrement++
		car.Id = c.autoIncrement
	}

	c.cars.Store(car.Id, car)
	return nil
}

func (c *repositoryInmem) Find(_ context.Context, id int) (*Car, error) {
	carRaw, ok := c.cars.Load(id)
	if !ok {
		return nil, ErrNotFound
	}

	return carRaw.(*Car), nil
}

func (c *repositoryInmem) FindAll(_ context.Context) ([]*Car, error) {
	allCars := make([]*Car, 0)

	c.cars.Range(func(key, value interface{}) bool {
		allCars = append(allCars, value.(*Car))

		return true
	})

	return allCars, nil
}

func (c *repositoryInmem) Del(_ context.Context, id int) error {
	c.cars.Delete(id)
	return nil
}
