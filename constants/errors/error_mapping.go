package errors

import (
	errorOrder "order-services/constants/errors/order"
)

func ErrMapping(err error) bool {
	var (
		GeneralErrors = GeneralErrors
		OrderErrors   = errorOrder.OrderErrors
	)

	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, OrderErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
