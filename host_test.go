package inventory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostMarshalBinary(t *testing.T) {
	host := Host{
		Name:     "a",
		Hostname: "a",
		IP:       "10.0.0.1",
		Roles:    []string{"a"},
	}
	expected, _ := json.Marshal(host)

	actual, err := host.MarshalBinary()

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestHostUnmarshalBinary(t *testing.T) {
	actual := Host{}
	expected := Host{
		Name:     "a",
		Hostname: "a",
		IP:       "10.0.0.1",
		Roles:    []string{"a"},
	}
	bytes, _ := json.Marshal(expected)

	err := actual.UnmarshalBinary(bytes)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}
