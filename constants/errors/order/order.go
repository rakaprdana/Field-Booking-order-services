package error

import "errors"

var (
	ErrOrderNotFound      = errors.New("Order not found")
	ErrFieldAlreadyBooked = errors.New("Field schedule already booked")
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrFieldAlreadyBooked,
}
