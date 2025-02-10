package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeEscapedJSONString(t *testing.T) {
	sampleString := "\"113c3a32a9d320b72190a04a020a0db3976ef36972673258e9a38a364f3dc3b0\""

	result, err := DecodeEscapedJSONString(sampleString)
	assert.NoError(t, err)
	assert.Equal(t, result, "113c3a32a9d320b72190a04a020a0db3976ef36972673258e9a38a364f3dc3b0")
}
