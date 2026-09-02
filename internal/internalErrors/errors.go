package internalErrors

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	JsonTypeError          = errors.New("Проверьте передаваемые типы данных")
	EmptyFieldsErr         = errors.New("Пожалуйста, не оставляйте пустые поля")
	NegativePriceError     = errors.New("Стоимость не может быть отрицательной")
	DataNotFound           = errors.New("Данные не найдены")
	IdRequestError         = errors.New("Проверьте передаваемый id")
	ParametresRequestError = errors.New("Проверьте передаваемые параметры")
	PriceError             = errors.New("Цена ОТ не может быть больше цены ДО")
)

func HandleDbError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return DataNotFound
	}
	return err
}
