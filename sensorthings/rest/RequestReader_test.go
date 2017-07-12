package rest

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http/httptest"
)

func TestReadRefFromWildcard(t *testing.T) {
	// arrange
	router := mux.NewRouter()
	router.HandleFunc("/v1.0/Things{id}/{params}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qo, _ := getQueryOptions(r)
		w.Write([]byte(fmt.Sprintf("%v", *qo.Ref)))
	}))

	ts := httptest.NewServer(router)
	defer ts.Close()

	// act
	resp, _ := http.Get(ts.URL + "/v1.0/Things(35)/$ref")

	// assert
	assert.True(t, resp != nil)
	assert.True(t, http.StatusOK == resp.StatusCode)
	body := resp.Body
	result, _ := ioutil.ReadAll(body)
	assert.True(t, string(result) == "true")
}

func TestReturnNoQueryOptionsOnFailedParse(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things?$count=none", nil)

	// act
	qo, err := getQueryOptions(req)

	// assert
	assert.Nil(t, qo)
	assert.NotNil(t, err)
}

func TestGetEntityId(t *testing.T) {
	// arrange
	router := mux.NewRouter()
	router.HandleFunc("/v1.0/Things{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := getEntityID(r)
		w.Write([]byte(id))
		// fmt.Println("func called")
	}))

	ts := httptest.NewServer(router)
	defer ts.Close()

	// act
	resp, _ := http.Get(ts.URL + "/v1.0/Things(35)")

	// assert
	assert.True(t, resp != nil)
	assert.True(t, http.StatusOK == resp.StatusCode)
	body := resp.Body
	result, _ := ioutil.ReadAll(body)
	assert.True(t, string(result) == "35")
}

func TestGetQueryOptions(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things/$value?$top=201", nil)

	// act
	qo, _ := getQueryOptions(req)

	// assert
	assert.True(t, qo != nil)
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
