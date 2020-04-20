package service

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
	car, err := s.CarsRepository.Find(ctx, id)
	return car, err
}

func (s *ServiceImpl) PostCar(ctx context.Context, car *Car) error {
	err := s.CarsRepository.Store(ctx, car)
	return err
}

func (s *ServiceImpl) PutCar(ctx context.Context, id int, carUpdated *Car) error {
	_, err := s.CarsRepository.Find(ctx, id)
	if err != nil {
		return err
	}

	carUpdated.Id = id

	err = s.CarsRepository.Store(ctx, carUpdated)
	return err
}

func (s *ServiceImpl) DeleteCar(ctx context.Context, id int) error {
	err := s.CarsRepository.Del(ctx, id)
	return err
}
