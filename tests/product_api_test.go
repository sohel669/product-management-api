package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"Project/api/routers"
)

func TestCreateProduct(t *testing.T) {
	router := routers.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
