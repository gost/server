package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTopics(t *testing.T) {
	// arrange
	// act
	topics := CreateTopics("GOST")
	// assert
	assert.True(t, len(topics) > 0, "Must have more than zero topics")
}
