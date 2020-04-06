package car

import "context"

// Service is a simple CRUD interface for cars.
type Service interface {
	PostCar(ctx context.Context, car *Car) error
	GetCar(ctx context.Context, id int) (*Car, error)
	PutCar(ctx context.Context, id int, car *Car) error
	PatchCar(ctx context.Context, id int, car *Car) error
	DeleteCar(ctx context.Context, id int) error
}

type service struct {
	cars Repository
}

func (s *service) GetCar(ctx context.Context, id int) (*Car, error) {
	car, err := s.GetCar(ctx, id)

	return car, err
}

func (s *service) PutCar(ctx context.Context, id int, car *Car) error {
	err := s.cars.Store(car)
	return err
}

func (s *service) PatchCar(ctx context.Context, id int, car *Car) error {
	err := s.cars.Store(car)
	return err
}

func (s *service) DeleteCar(ctx context.Context, id int) error {
	err := s.cars.Del(id)
	return err
}



