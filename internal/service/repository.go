package service

import (
	"context"
)

type Repository interface {
	Store(ctx context.Context, car *Car) error
	Find(ctx context.Context, id int) (*Car, error)
	FindAll(ctx context.Context) ([]*Car, error)
	Del(ctx context.Context, id int) error
}
