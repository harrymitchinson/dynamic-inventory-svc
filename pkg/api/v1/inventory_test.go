package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
	"github.com/harrymitchinson/dynamic-inventory-svc/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetInventory(t *testing.T) {

	t.Run("Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		svc := &mocks.HostService{}
		svc.On("GetHosts", env).Return([]inventory.Host{}, errors.New("failed"))

		NewInventoryController(svc).GetInventory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		svc := &mocks.HostService{}
		hosts := []inventory.Host{
			{
				Name: "a", Hostname: "a",
				IP:    "10.0.0.1",
				Roles: []string{"a"},
			},
		}
		svc.On("GetHosts", env).Return(hosts, nil)

		NewInventoryController(svc).GetInventory(c)

		if assert.Equal(t, http.StatusOK, w.Code) {
			expectedBytes, _ := json.Marshal(inventory.NewInventoryFile(hosts))
			actualBytes, _ := ioutil.ReadAll(w.Body)
			expectedStr := string(expectedBytes)
			actualStr := string(actualBytes)
			assert.JSONEq(t, expectedStr, actualStr)
		}
	})
}
