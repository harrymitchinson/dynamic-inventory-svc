package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnvironment(t *testing.T) {
	cases := map[string]Environment{
		"dev": EnvironmentDev,
		"dEv": EnvironmentDev,
		"DEV": EnvironmentDev,
		"stg": EnvironmentStg,
		"Stg": EnvironmentStg,
		"STG": EnvironmentStg,
		"prd": EnvironmentPrd,
		"prD": EnvironmentPrd,
		"PRD": EnvironmentPrd,
	}

	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			actual, err := ParseEnvironment(input)
			if assert.NoError(t, err) {
				assert.Equal(t, expected, actual)
			}
		})
	}

	t.Run("invalid", func(t *testing.T) {
		_, err := ParseEnvironment("invalid")
		if assert.Error(t, err) {
			assert.Equal(t, ErrInvalidEnvironment, err)
		}
	})
}
