package odata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryValidator(t *testing.T) {
	// arrange
	assert.Equal(t, true, IsValidOdataQuery("$filter=name eq 'ho'"))
	assert.Equal(t, false, IsValidOdataQuery("$notexisting=name eq 'ho'"))
}
