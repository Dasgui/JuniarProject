package internalErrors

import "errors"

var (
	JsonTypeError      = errors.New("Проверьте передаваемые типы данных")
	EmptyFieldsErr     = errors.New("Пожалуйста, не оставляйте пустые поля")
	NegativePriceError = errors.New("Стоимость не может быть отрицательной")
)
