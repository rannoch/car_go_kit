package car

import "time"

type Car struct {
	Id      int       `json:"id"`
	Brand   string    `json:"brand"`
	Model   string    `json:"model"`
	Created time.Time `json:"created"`
}
