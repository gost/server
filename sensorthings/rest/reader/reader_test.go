package reader

import (
	"net/http"
	"testing"

	entities "github.com/gost/core"
	gostErrors "github.com/gost/server/errors"
	"github.com/stretchr/testify/assert"

	"bytes"
	"io/ioutil"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func TestGetEntityId(t *testing.T) {
	// arrange
	router := mux.NewRouter()
	router.HandleFunc("/v1.0/Things{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetEntityID(r)
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

func TestCheckContentTypeWithoutHeadersShouldReturnFalse(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things(1)", nil)
	w := httptest.NewRecorder()

	// act
	res := CheckContentType(w, req, false)

	// assert
	assert.True(t, res)
}

func TestCheckContentTypeWithContentTypeHeaderShouldReturnTrue(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things(1)", nil)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// act
	res := CheckContentType(w, req, false)

	// assert
	assert.True(t, res)
}

func TestCheckContentTypeWithoutContentTypeHeaderShouldReturnFalse(t *testing.T) {
	// arrange
	req, _ := http.NewRequest("GET", "/v1.0/Things(1)", nil)
	req.Header.Add("Content-Type", "superformat")
	w := httptest.NewRecorder()

	// act
	res := CheckContentType(w, req, false)

	// assert
	assert.False(t, res)
}

func TestCheckAndGetBodyWithNoBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bla", nil)

	// act
	CheckAndGetBody(rr, req, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCheckAndGetBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bla", bytes.NewReader([]byte("")))

	// act
	CheckAndGetBody(rr, req, false)

	// assert
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestParseEntity(t *testing.T) {
	// arrange
	thing := &entities.Thing{}
	thingBytes := []byte("{\"name\": \"thing1\", \"description\": \"test thing 1\"}")

	location := &entities.Location{}
	historicalLocation := &entities.HistoricalLocation{}
	datastream := &entities.Datastream{}
	sensor := &entities.Sensor{}
	observedProperty := &entities.ObservedProperty{}
	observation := &entities.Observation{}
	featureOfinterest := &entities.FeatureOfInterest{}

	// act
	tErr := ParseEntity(thing, thingBytes)
	lErr := ParseEntity(location, nil)
	hlErr := ParseEntity(historicalLocation, nil)
	dErr := ParseEntity(datastream, nil)
	sErr := ParseEntity(sensor, nil)
	opErr := ParseEntity(observedProperty, nil)
	oErr := ParseEntity(observation, nil)
	fErr := ParseEntity(featureOfinterest, nil)

	// assert
	assert.Nil(t, tErr)
	assert.Equal(t, thing.Name, "thing1")
	assert.Equal(t, thing.Description, "test thing 1")
	assert.Equal(t, 400, getStatusCode(lErr))
	assert.Equal(t, 400, getStatusCode(hlErr))
	assert.Equal(t, 400, getStatusCode(dErr))
	assert.Equal(t, 400, getStatusCode(sErr))
	assert.Equal(t, 400, getStatusCode(opErr))
	assert.Equal(t, 400, getStatusCode(oErr))
	assert.Equal(t, 400, getStatusCode(fErr))
}

func getStatusCode(err error) int {
	switch e := err.(type) {
	case gostErrors.APIError:
		return e.GetHTTPErrorStatusCode()
	}

	return 0
}
