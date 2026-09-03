package test_data

import (
	"JuniarProject/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestTable struct {
	Name        string
	QueryParams string

	Setup func()

	Except         []models.Product
	ExpectedStatus int
	ExpectedError  error
	ExpectedCount  int
}

func StartTest(t *testing.T, testTable []TestTable, handler func(w http.ResponseWriter, r *http.Request)) {
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Setup()

			req := httptest.NewRequest(http.MethodGet, "/products"+tt.QueryParams, nil)
			w := httptest.NewRecorder()

			// h.GetProducts(w, req)
			handler(w, req)

			assert.Equal(t, tt.ExpectedStatus, w.Code)

			if tt.ExpectedStatus == http.StatusOK {
				var products []models.Product
				err := json.Unmarshal(w.Body.Bytes(), &products)
				require.NoError(t, err)
				assert.Len(t, products, tt.ExpectedCount)
			} else if tt.ExpectedError != nil {
				assert.Contains(t, strings.TrimSpace(w.Body.String()), tt.ExpectedError.Error())
			}
		})
	}
}
