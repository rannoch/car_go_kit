package car

import (
	"sync"
)

type repositoryInmem struct {
	cars          sync.Map
	autoIncrement int
}

func NewCarRepositoryInmem() *repositoryInmem {
	return &repositoryInmem{}
}

func (c *repositoryInmem) Store(car *Car) error {
	if car.Id == 0 {
		c.autoIncrement++
		car.Id = c.autoIncrement
	}

	c.cars.Store(car.Id, car)
	return nil
}

func (c *repositoryInmem) Find(id int) (*Car, error) {
	carRaw, ok := c.cars.Load(id)
	if !ok {
		return nil, ErrNotFound
	}

	return carRaw.(*Car), nil
}

func (c *repositoryInmem) FindAll() ([]*Car, error) {
	allCars := make([]*Car, 0)

	c.cars.Range(func(key, value interface{}) bool {
		allCars = append(allCars, value.(*Car))

		return true
	})

	return allCars, nil
}

func (c *repositoryInmem) Del(id int) error {
	c.cars.Delete(id)
	return nil
}
