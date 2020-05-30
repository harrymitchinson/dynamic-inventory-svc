package inventory

import (
	"errors"
	"strings"
)

// Environment is the list of allowed values for the inventory environment.
type Environment string

// List of values that Environment can take.
const (
	EnvironmentDev Environment = "dev"
	EnvironmentStg Environment = "stg"
	EnvironmentPrd Environment = "prd"
)

// List of errors that the Environment can return.
var (
	ErrInvalidEnvironment = errors.New("invalid environment")
)

// ParseEnvironment parses a string to an Environment
func ParseEnvironment(v string) (Environment, error) {
	e := Environment(strings.ToLower(v))

	switch e {
	case EnvironmentDev:
		return EnvironmentDev, nil
	case EnvironmentStg:
		return EnvironmentStg, nil
	case EnvironmentPrd:
		return EnvironmentPrd, nil
	default:
		return Environment(""), ErrInvalidEnvironment
	}
}
