package test_data

import (
	"JuniarProject/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ProductTestTable struct {
	Name        string
	QueryParams string
	Body        string

	Setup func()

	Except         interface{}
	ExpectedStatus int
	ExpectedError  error
	ExpectedCount  int
}

func StartTest(t *testing.T, testTable []ProductTestTable, r *chi.Mux, method string) {
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Setup()

			req := httptest.NewRequest(method, "/products/"+tt.QueryParams, strings.NewReader(tt.Body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedStatus, w.Code)

			if tt.ExpectedStatus == http.StatusOK {
				if tt.ExpectedCount != 0 {
					var products []models.Product
					err := json.Unmarshal(w.Body.Bytes(), &products)
					require.NoError(t, err)
					assert.Len(t, products, tt.ExpectedCount)
				} else {
					var product models.Product
					err := json.Unmarshal(w.Body.Bytes(), &product)
					require.NoError(t, err)
					assert.Equal(t, tt.Except, product)
				}
			} else if tt.ExpectedError != nil {
				assert.Contains(t, strings.TrimSpace(w.Body.String()), tt.ExpectedError.Error())
			}
		})
	}
}
