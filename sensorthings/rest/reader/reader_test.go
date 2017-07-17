package reader

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http/httptest"
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

/*func TestCheckAndGetBodyWithWrongBody(t *testing.T) {
	// arrange
	rr := httptest.NewRecorder()
	var reader io.Reader
	b, _ := json.Marshal("")
	reader = bytes.NewReader(b)
	req, _ := http.NewRequest("GET", "/bla", reader)

	// act
	CheckAndGetBody(rr, req, false)

	// assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}*/
