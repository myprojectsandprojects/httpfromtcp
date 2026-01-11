package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParse(t *testing.T) {
	// Valid single header
	headers := Create()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NotNil(t, headers) // ?
	assert.Equal(t, 23, n)
	assert.False(t, done)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])

	// Allowed special characters
	headers = Create()
	data = []byte("!#$%&'*+-.^_`|~: anything is allowd in a field-value\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 54, n)
	assert.False(t, done)
	require.NoError(t, err)
	assert.Equal(t, "anything is allowd in a field-value", headers["!#$%&'*+-.^_`|~"])

	// Disallowed special character
	headers = Create()
	data = []byte("@!#$%&'*+-.^_`|~: anything is allowd in a field-value\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	require.Error(t, err)

	// Valid spacing in field-name
	headers = Create()
	data = []byte(" Abc: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 30, n)
	assert.False(t, done)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["abc"])

	// Invalid spacing in field-name
	headers = Create()
	data = []byte(" Abc : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	require.Error(t, err)

	// This is when we should be done
	headers = Create()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 2, n)
	assert.True(t, done)
	require.NoError(t, err)

	// Invalid characters in field-name
	headers = Create()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	require.Error(t, err)

	// Field-name can appear multiple times in headers
	headers = Create()
	headers["set-person"] = "prime-loves-zig"
	data = []byte("Set-Person: lane-loves-go\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 27, n)
	assert.False(t, done)
	require.NoError(t, err)
	assert.Equal(t, "prime-loves-zig, lane-loves-go", headers["set-person"])
	assert.Equal(t, "", headers["Set-Person"])
}
