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

var (
	JsonTypeError         = errors.New("Проверьте передаваемые типы данных")
	JsonUnknownFieldError = errors.New("Присутствует неизвестное поле")
	JsonEmptyBodyError    = errors.New("Пустое тело запроса")
	JsonSyntaxError       = errors.New("Некорректный JSON формат")

	EmptyFieldsError      = errors.New("Не все поля заполнены")
	NegativePriceError    = errors.New("Стоимость не может быть отрицательной")
	DataNotFound          = errors.New("Данные не найдены")
	InvalidParameterError = errors.New("Проверьте передаваемые параметры")
	IdRangeError          = errors.New("Id выходит за допустимые пределы")
	PriceError            = errors.New("Цена ОТ не может быть больше цены ДО")
	InternalServerErr     = errors.New("Ошибка в работе сервера")
)

func HandleDbError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return DataNotFound
	}
	return err
}

func ParseUrlError(err error) error {
	if err == nil {
		return nil
	}

	// Проверяем на синтаксическую ошибку (не число)
	if errors.Is(err, strconv.ErrSyntax) {
		return InvalidParameterError
	}

	// Проверяем на выход за пределы
	if errors.Is(err, strconv.ErrRange) {
		return IdRangeError
	}

	return err
}

func ParseJsonError(err error) error {
	if err == nil {
		return nil
	}

	if err.Error() == "EOF" || err.Error() == "request body is empty" {
		return JsonEmptyBodyError
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return JsonSyntaxError
	}

	// Проверяем на ошибки типов
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		return JsonTypeError
	}

	if strings.Contains(err.Error(), "unknown field") {
		return JsonUnknownFieldError
	}

	return err
}

func PrintError(w http.ResponseWriter, err error) {
	msg, code := DecipherError(err)
	log.Println(err)
	http.Error(w, msg, code)
}

func DecipherError(err error) (string, int) {
	switch {
	case errors.Is(err, PriceError):
		return PriceError.Error(), http.StatusBadRequest

	case errors.Is(err, IdRangeError):
		return IdRangeError.Error(), http.StatusBadRequest

	case errors.Is(err, InvalidParameterError):
		return InvalidParameterError.Error(), http.StatusBadRequest

	case errors.Is(err, JsonSyntaxError):
		return JsonSyntaxError.Error(), http.StatusBadRequest

	case errors.Is(err, JsonEmptyBodyError):
		return JsonEmptyBodyError.Error(), http.StatusBadRequest

	case errors.Is(err, JsonUnknownFieldError):
		return JsonUnknownFieldError.Error(), http.StatusBadRequest

	case errors.Is(err, JsonTypeError):
		return JsonTypeError.Error(), http.StatusBadRequest
	case errors.Is(err, EmptyFieldsError):
		return EmptyFieldsError.Error(), http.StatusBadRequest

	case errors.Is(err, NegativePriceError):
		return NegativePriceError.Error(), http.StatusBadRequest

	case errors.Is(err, DataNotFound):
		return DataNotFound.Error(), http.StatusNotFound
	default:
		return InternalServerErr.Error(), http.StatusInternalServerError
	}
}
