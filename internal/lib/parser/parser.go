package parser

import (
	"JuniarProject/internal/internalErrors"
	"net/http"
	"strconv"
)

func ParseRequestToFloat(r *http.Request, key string) (float64, error) {
	if priceFromStr := r.URL.Query().Get(key); priceFromStr != "" {
		if val, err := strconv.ParseFloat(priceFromStr, 64); err == nil {
			if val < 0 {
				return -1.0, internalErrors.NegativePriceError.Err
			}
			return val, nil
		} else {
			return -1.0, err
		}
	}

	return 0.0, nil
}

func ParseRequestToInt(r *http.Request, key string) (int, error) {
	if limitStr := r.URL.Query().Get(key); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err != nil || val < 0 {
			return -1, err
		} else {
			return val, nil
		}
	}

	return 0, nil
}
