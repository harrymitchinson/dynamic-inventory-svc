package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestEnvironmentMiddleware(t *testing.T) {
	t.Run("NoEnvironment", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		EnvironmentMiddleware(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		_, exists := c.Get(EnvironmentKey)
		assert.False(t, exists)
	})

	t.Run("InvalidEnvironment", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: EnvironmentKey, Value: "invalid"}}

		EnvironmentMiddleware(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		_, exists := c.Get(EnvironmentKey)
		assert.False(t, exists)
	})

	t.Run("ValidEnvironment", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		expected := inventory.EnvironmentDev
		c.Params = []gin.Param{{Key: EnvironmentKey, Value: string(expected)}}

		EnvironmentMiddleware(c)

		assert.Equal(t, http.StatusOK, w.Code)
		actual := c.MustGet(EnvironmentKey).(inventory.Environment)
		assert.Equal(t, actual, expected)
	})
}
