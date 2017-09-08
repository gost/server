package postgis

import (
	"testing"

	"github.com/gost/godata"
	"github.com/stretchr/testify/assert"
)

func init() {
	setupLogger()
}

func TestToIntIDForString(t *testing.T) {
	// arrange
	fid := "4"

	// act
	intID, err := ToIntID(fid)

	// assert
	assert.True(t, intID == 4)
	assert.True(t, err)
}

func TestToIntIDForFloat(t *testing.T) {
	// arrange
	fid := 6.4

	// act
	intID, err := ToIntID(fid)

	// assert
	assert.True(t, intID == 6)
	assert.True(t, err)
}

func TestContainsToLower(t *testing.T) {
	// arrange
	ss, _ := godata.ParseSelectString("Hallo")
	search := "HALLO"

	// act
	res := ContainsToLower(ss.SelectItems, search)

	// assert
	assert.True(t, res)
}

func TestContainsNotToLower(t *testing.T) {
	// arrange
	ss, _ := godata.ParseSelectString("Halllo")
	search := "HALLO"

	// act
	res := ContainsToLower(ss.SelectItems, search)

	// assert
	assert.False(t, res)
}

func TestJsonToMapSucceeds(t *testing.T) {
	// arrange
	jsonstring := `{"value": [{"name": "Things","url": "http://gost.geodan.nl/v1.0/Things"}]}`

	// act
	res, err := JSONToMap(&jsonstring)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res["value"])
}

func TestJsonToMapFails(t *testing.T) {
	// arrange
	jsonstring := ``

	// act
	_, err := JSONToMap(&jsonstring)

	// assert
	assert.Nil(t, err)
}

func TestJsonToMapFailsWithWrongData(t *testing.T) {
	// arrange
	jsonstring := `hoho`

	// act
	_, err := JSONToMap(&jsonstring)

	// assert
	assert.NotNil(t, err)
}
