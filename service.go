package car

import (
	"context"
	"errors"
)

// Service is a simple CRUD interface for CarsRepository.
type Service interface {
	GetCar(ctx context.Context, id int) (*Car, error)
	PostCar(ctx context.Context, car *Car) error
	PutCar(ctx context.Context, id int, car *Car) error
	DeleteCar(ctx context.Context, id int) error
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type ServiceImpl struct {
	CarsRepository Repository
}

func (s *ServiceImpl) GetCar(ctx context.Context, id int) (*Car, error) {
	car, err := s.CarsRepository.Find(id)
	return car, err
}

func (s *ServiceImpl) PostCar(ctx context.Context, car *Car) error {
	if car.Id != 0 {
		return ErrBadRouting
	}

	err := s.CarsRepository.Store(car)
	return err
}

func (s *ServiceImpl) PutCar(ctx context.Context, id int, carUpdated *Car) error {
	car, err := s.CarsRepository.Find(id)
	if err != nil {
		return err
	}

	car = carUpdated
	car.Id = id

	err = s.CarsRepository.Store(car)
	return err
}

func (s *ServiceImpl) DeleteCar(ctx context.Context, id int) error {
	err := s.CarsRepository.Del(id)
	return err
}



