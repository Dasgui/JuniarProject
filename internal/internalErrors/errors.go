package internalErrors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type AppError struct {
	Err  error
	Code int
}

var (
	JsonTypeError = AppError{
		errors.New("Проверьте передаваемые типы данных"),
		http.StatusBadRequest,
	}
	JsonUnknownFieldError = AppError{
		errors.New("Присутствует неизвестное поле"),
		http.StatusBadRequest,
	}
	JsonEmptyBodyError = AppError{
		errors.New("Пустое тело запроса"),
		http.StatusBadRequest,
	}
	JsonSyntaxError = AppError{
		errors.New("Некорректный JSON формат"),
		http.StatusBadRequest,
	}

	EmptyFieldsError = AppError{
		errors.New("Не все поля заполнены"),
		http.StatusBadRequest,
	}
	NegativePriceError = AppError{
		errors.New("Стоимость не может быть отрицательной"),
		http.StatusBadRequest,
	}

	DataNotFound = AppError{
		errors.New("Данные не найдены"),
		http.StatusNotFound,
	}

	// InvalidParameterError - структура для ответа с ошибкой в Swagger
	InvalidParameterError = AppError{
		errors.New("Проверьте передаваемые параметры"),
		http.StatusBadRequest,
	}
	IdRangeError = AppError{
		errors.New("Id выходит за допустимые пределы"),
		http.StatusBadRequest,
	}

	PriceError = AppError{
		errors.New("Цена ОТ не может быть больше цены ДО"),
		http.StatusBadRequest,
	}
	InternalServerErr = AppError{
		errors.New("Ошибка в работе сервера"),
		http.StatusInternalServerError,
	}
)

func PrintError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	log.Println(err)
	appError := ConvertError(err)
	http.Error(w, appError.Err.Error(), appError.Code)
}

func ConvertError(err error) AppError {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError
	var unmarshalInvalidErr *json.InvalidUnmarshalError

	if err.Error() == "EOF" || err.Error() == "request body is empty" {
		return JsonEmptyBodyError
	}

	switch {
	case errors.As(err, &unmarshalInvalidErr):
		return JsonTypeError

	case strings.Contains(err.Error(), "unknown field"):
		return JsonUnknownFieldError

	case errors.As(err, &unmarshalTypeErr):
		return JsonTypeError

	case errors.As(err, &syntaxErr):
		return JsonSyntaxError
	}

	switch {
	case errors.Is(err, PriceError.Err):
		return PriceError
	case errors.Is(err, NegativePriceError.Err):
		return NegativePriceError
	case errors.Is(err, EmptyFieldsError.Err):
		return EmptyFieldsError
	case errors.Is(err, strconv.ErrRange):
		return IdRangeError
	case errors.Is(err, strconv.ErrSyntax):
		return InvalidParameterError
	case errors.Is(err, pgx.ErrNoRows):
		return DataNotFound
	}

	return InternalServerErr
}
