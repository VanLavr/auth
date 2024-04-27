package hasher_test

import (
	"testing"

	"github.com/VanLavr/auth/internal/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	data := "hello world"
	data2 := "hello"

	encoded := hasher.Hshr.Encrypt(data)

	assert := assert.New(t)
	assert.Equal(true, hasher.Hshr.Validate(encoded, data))
	assert.Equal(false, hasher.Hshr.Validate(encoded, data2))
}
