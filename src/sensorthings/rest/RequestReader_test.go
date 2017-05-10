package rest

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"net/http/httptest"
)


func TestGetQueryOptions(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things?$top=1", nil)

	// act
	qo, _ := getQueryOptions(req)

	// assert
	assert.True(t, qo != nil)
}

func TestGetQueryOptionsWithWrongOptions(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things?$hoho", nil)

	// act
	qo, _ := getQueryOptions(req)

	// assert
	assert.Nil(t, qo)
}

func TestCheckContentTypeWithoutHeadersShouldReturnFalse(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things(1)", nil)
	w := httptest.NewRecorder()

	// act
	res := checkContentType(w, req)

	// assert
	assert.True(t, res)
}

func TestCheckContentTypeWithContentTypeHeaderShouldReturnTrue(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things(1)", nil)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// act
	res := checkContentType(w, req)

	// assert
	assert.True(t, res)
}

func TestCheckContentTypeWithoutContentTypeHeaderShouldReturnFalse(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things(1)", nil)
	req.Header.Add("Content-Type", "superformat")
	w := httptest.NewRecorder()

	// act
	res := checkContentType(w, req)

	// assert
	assert.False(t, res)
}
