package writer

import (
	"errors"
	"github.com/gost/godata"
	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/odata"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gost/server/sensorthings/models"
	"fmt"
	"io/ioutil"
)

func TestSendErrorWithNoError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()

	// act
	SendError(rr, nil, true)

	// assert
	assert.True(t, rr.Code == http.StatusInternalServerError)
}

func TestSendErrorWithNoIdentJson(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.Name = "yo"

	// act
	SendJSONResponse(rr, http.StatusTeapot, thing, nil, false)

	// assert
	assert.True(t, rr.Code == http.StatusTeapot)
}

func TestSendErrorWithError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	err1 := errors.New("wrong")
	errs := []error{err1}

	// act
	SendError(rr, errs, false)

	// assert
	assert.True(t, rr.Code == http.StatusInternalServerError)
}

func TestSendJsonResponseWithNoData(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()

	// act
	SendJSONResponse(rr, http.StatusTeapot, nil, nil, false)

	// assert
	assert.True(t, rr.Code == http.StatusOK)
}

func TestSendJsonResponseWithData(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.Name = "yo"

	// act
	SendJSONResponse(rr, http.StatusTeapot, thing, nil, false)

	// assert
	assert.True(t, rr.Code == http.StatusTeapot)
}

func TestSendJsonResponseWithDataAndQueryOptions(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{}
	thing.Name = "yo"
	req, _ := http.NewRequest("GET", "/v1.0/Things?$top=1&$select=name,id,description", nil)
	qo, _ := odata.GetQueryOptions(req, 20)

	val := odata.GoDataValueQuery(true)
	qo.Value = &val

	// act
	SendJSONResponse(rr, http.StatusTeapot, thing, qo, false)

	// assert
	assert.True(t, rr.Code == http.StatusTeapot)
}

func TestSendJsonResponseErrorOnMarshalError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	c := make(chan int)
	m := map[string]interface{}{"chan": c}

	// assert
	assert.Panics(t, func() { SendJSONResponse(rr, http.StatusTeapot, m, nil, false) })
}

func TestSendJsonResponseWithRefAndNovalueError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{Name: "yo"}
	qo := &odata.QueryOptions{}
	valQuery := odata.GoDataValueQuery(true)
	qo.Value = &valQuery
	qo.Select = &godata.GoDataSelectQuery{SelectItems: nil}

	// act
	SendJSONResponse(rr, http.StatusTeapot, thing, qo, false)

	// assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRequestValueWithNonexistingParam(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	thing := &entities.Thing{Name: "yo"}
	qo := &odata.QueryOptions{}
	valQuery := odata.GoDataValueQuery(true)
	qo.Value = &valQuery
	qo.Select = &godata.GoDataSelectQuery{SelectItems: []*godata.SelectItem{{Segments: []*godata.Token{{Value: "nonexistingparam"}}}}}

	// act
	SendJSONResponse(rr, http.StatusTeapot, thing, qo, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEncodingNotSupported(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	err := []error{errors.New("Encoding not supported")}

	// act
	SendError(rr, err, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCountCollection(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	qo := &odata.QueryOptions{}
	c := odata.GoDataCollectionCountQuery(true)
	qo.CollectionCount = &c
	ar := models.ArrayResponse{ Count: 10 }

	// act
	SendJSONResponse(rr, http.StatusOK, ar, qo, false)
	body, _ := ioutil.ReadAll(rr.Body)

	// assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(body), fmt.Sprintf("%v", ar.Count))
}

func TestCountCollectionError(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	qo := &odata.QueryOptions{}
	c := odata.GoDataCollectionCountQuery(true)
	qo.CollectionCount = &c
	ar := models.ArrayResponse{}

	// act
	SendJSONResponse(rr, http.StatusOK, ar, qo, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}