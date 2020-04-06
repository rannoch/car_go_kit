package car

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServiceImpl_GetCar(t *testing.T) {
	var carService ServiceImpl
	carService.CarsRepository = NewCarRepositoryInmem()

	var car *Car
	var err error

	car, err = carService.GetCar(nil, 1)

	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, car)

	car = &Car{
		Brand:   "Toyota",
		Model:   "Corolla",
		Created: time.Now(),
	}

	err = carService.PostCar(nil, car)

	assert.NoError(t, err)

	c, err := carService.GetCar(nil, car.Id)

	assert.NoError(t, err)
	assert.Equal(t, car, c)
}

func TestServiceImpl_PostCar(t *testing.T) {
	var carService ServiceImpl
	carService.CarsRepository = NewCarRepositoryInmem()

	var car *Car
	var err error

	car = &Car{
		Brand:   "Toyota",
		Model:   "Corolla",
		Created: time.Now(),
	}

	err = carService.PostCar(nil, car)

	assert.NoError(t, err)
}

func TestServiceImpl_PutCar(t *testing.T) {
	var carService ServiceImpl
	carService.CarsRepository = NewCarRepositoryInmem()

	var car *Car
	var err error

	car = &Car{
		Brand:   "Toyota",
		Model:   "Corolla",
		Created: time.Now(),
	}

	err = carService.PutCar(nil, car.Id, car)

	assert.Equal(t, ErrNotFound, err)

	err = carService.PostCar(nil, car)

	assert.NoError(t, err)

	car.Brand = "Mercedes"
	car.Model = "Benz"

	err = carService.PutCar(nil, car.Id, car)

	assert.NoError(t, err)

	c, err := carService.GetCar(nil, car.Id)

	assert.NoError(t, err)
	assert.Equal(t, car, c)
}

func TestServiceImpl_DeleteCar(t *testing.T) {
	var carService ServiceImpl
	carService.CarsRepository = NewCarRepositoryInmem()

	var car *Car
	var err error

	car = &Car{
		Brand:   "Toyota",
		Model:   "Corolla",
		Created: time.Now(),
	}

	err = carService.PostCar(nil, car)

	assert.NoError(t, err)

	err = carService.DeleteCar(nil, car.Id)

	assert.NoError(t, err)

	car, err = carService.GetCar(nil, car.Id)

	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, car)
}
