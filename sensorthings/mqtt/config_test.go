package mqtt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTopics(t *testing.T) {
	// arrange
	// act
	topics := CreateTopics()
	// assert
	assert.True(t, len(topics) > 0, "Must have more than zero topics")
}
