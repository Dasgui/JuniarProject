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
	//msg, code := DecipherError(err)
	log.Println(err)
	appError := ConvertError(err)
	http.Error(w, appError.Err.Error(), appError.Code)
}

func ConvertError(err error) AppError {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError
	var unmarshalInvalidErr *json.InvalidUnmarshalError

	//var appErr AppError
	//if errors.As(err, &appErr) {
	//	return appErr
	//}

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
	//if errors.As(err, &syntaxErr) {
	//	return JsonSyntaxError
	//}

	//if errors.As(err, &unmarshalTypeErr) {
	//	return JsonTypeError
	//}

	//if strings.Contains(err.Error(), "unknown field") {
	//	return JsonUnknownFieldError
	//}

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

//func HandleDbError(err error) AppError {
//	if errors.Is(err, pgx.ErrNoRows) {
//		return DataNotFound
//	}
//	return InternalServerErr
//}

//func ParseUrlError(err error) error {
//if err == nil {
//	return nil
//}

// Проверяем на синтаксическую ошибку (не число)
//if errors.Is(err, strconv.ErrSyntax) {
//	return InvalidParameterError
//}

// Проверяем на выход за пределы
//if errors.Is(err, strconv.ErrRange) {
//	return IdRangeError
//}
//
//	return err
//}

//func ParseJsonError(err error) error {
//	if err == nil {
//		return nil
//	}

//if err.Error() == "EOF" || err.Error() == "request body is empty" {
//	return JsonEmptyBodyError
//}
//
//var syntaxErr *json.SyntaxError
//if errors.As(err, &syntaxErr) {
//	return JsonSyntaxError
//}
//
//// Проверяем на ошибки типов
//var unmarshalTypeErr *json.UnmarshalTypeError
//if errors.As(err, &unmarshalTypeErr) {
//	return JsonTypeError
//}
//
//if strings.Contains(err.Error(), "unknown field") {
//	return JsonUnknownFieldError
//}
//
//	return err
//}

//func DecipherError(err error) (string, int) {
//	switch {
//case errors.Is(err, PriceError):
//	return PriceError.Error(), http.StatusBadRequest
//
//case errors.Is(err, IdRangeError):
//	return IdRangeError.Error(), http.StatusBadRequest
//
//case errors.Is(err, InvalidParameterError):
//	return InvalidParameterError.Error(), http.StatusBadRequest
//
//case errors.Is(err, JsonSyntaxError):
//	return JsonSyntaxError.Error(), http.StatusBadRequest
//
//case errors.Is(err, JsonEmptyBodyError):
//	return JsonEmptyBodyError.Error(), http.StatusBadRequest
//
//case errors.Is(err, JsonUnknownFieldError):
//	return JsonUnknownFieldError.Error(), http.StatusBadRequest
//
//case errors.Is(err, JsonTypeError):
//	return JsonTypeError.Error(), http.StatusBadRequest
//case errors.Is(err, EmptyFieldsError):
//	return EmptyFieldsError.Error(), http.StatusBadRequest
//
//case errors.Is(err, NegativePriceError):
//	return NegativePriceError.Error(), http.StatusBadRequest
//default:
//	return "", http.StatusInternalServerError
//}
//}
