package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
)

const (
	// EnvironmentKey is the key which the 'environment' parameter will be loaded from.
	EnvironmentKey = "environment"
)

var (
	// Validator is a singleton validator.Validate instance for the API to use for validating objects.
	Validator *validator.Validate = validator.New()
)

func getEnvironment(c *gin.Context) inventory.Environment {
	return c.MustGet(EnvironmentKey).(inventory.Environment)
}

// EnvironmentMiddleware validates the EnvrionmentKey param on the request, if the param is missing or invalid the response will be a BadRequest with ErrInvalidEnvironment. If the param is valid this is saved to the request context for later use within controllers.
func EnvironmentMiddleware(c *gin.Context) {
	env, err := inventory.ParseEnvironment(c.Param(EnvironmentKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.Set(EnvironmentKey, env)
	c.Next()
}
