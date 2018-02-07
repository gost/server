package odata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryValidator(t *testing.T) {
	// arrange
	assert.Equal(t, true, IsValidOdataQuery("$filter=name eq 'ho'"))
	assert.Equal(t, true, IsValidOdataQuery(fmt.Sprintf("%sfilter=name eq 'ho'", "%24")))
	assert.Equal(t, false, IsValidOdataQuery("$notexisting=name eq 'ho'"))
}
