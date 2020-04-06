package inmem

import (
	"github.com/rannoch/car"
	"sync"
)

type carRepositoryInmem struct {
	cars sync.Map
}

func NewCarRepositoryInmem() *carRepositoryInmem {
	return &carRepositoryInmem{}
}

func (c *carRepositoryInmem) Store(car *car.Car) error {
	c.cars.Store(car.Id, car)
	return nil
}

func (c *carRepositoryInmem) Find(id int) (*car.Car, error) {
	carRaw, ok := c.cars.Load(id)
	if !ok {
		return nil, car.ErrNotFound
	}

	return carRaw.(*car.Car), nil
}

func (c *carRepositoryInmem) FindAll() ([]*car.Car, error) {
	allCars := make([]*car.Car, 0)

	c.cars.Range(func(key, value interface{}) bool {
		allCars = append(allCars, value.(*car.Car))

		return true
	})

	return allCars, nil
}

func (c *carRepositoryInmem) Del(id int) error {
	c.cars.Delete(id)
	return nil
}
