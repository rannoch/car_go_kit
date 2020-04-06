package car

type Repository interface {
	Store(car *Car) error
	Find(id int) (*Car, error)
	FindAll() ([]*Car, error)
	Del(id int) error
}
