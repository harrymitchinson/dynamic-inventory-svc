package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
	"github.com/harrymitchinson/dynamic-inventory-svc/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetHosts(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		svc := &mocks.HostService{}
		svc.On("GetHosts", env).Return([]inventory.Host{}, errors.New("failed"))

		NewHostController(svc).GetHosts(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		hosts := []inventory.Host{
			{
				Name: "a", Hostname: "a",
				IP:    "10.0.0.1",
				Roles: []string{"a"},
			},
		}

		svc := &mocks.HostService{}
		svc.On("GetHosts", env).Return(hosts, nil)

		NewHostController(svc).GetHosts(c)

		if assert.Equal(t, http.StatusOK, w.Code) {
			expected, _ := json.Marshal(hosts)
			actual, _ := ioutil.ReadAll(w.Body)
			assert.Equal(t, expected, actual)
		}
	})

}

func TestCreateHost(t *testing.T) {
	t.Run("InvalidFormat", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		c.Request, _ = http.NewRequest("", "", strings.NewReader("not json"))

		NewHostController(&mocks.HostService{}).CreateHost(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		body, _ := json.Marshal(inventory.Host{})
		c.Request, _ = http.NewRequest("", "", bytes.NewReader(body))

		NewHostController(&mocks.HostService{}).CreateHost(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		host := &inventory.Host{
			Name:     "testing",
			Hostname: "testing",
			IP:       "10.0.0.1",
			Roles:    []string{"a"},
		}
		body, _ := json.Marshal(host)
		c.Request, _ = http.NewRequest("", "", bytes.NewReader(body))

		svc := &mocks.HostService{}
		svc.On("SetHost", env, host).Return(errors.New("failed"))

		NewHostController(svc).CreateHost(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		env := inventory.EnvironmentDev
		c.Keys = map[string]interface{}{EnvironmentKey: env}

		host := &inventory.Host{
			Name:     "testing",
			Hostname: "testing",
			IP:       "10.0.0.1",
			Roles:    []string{"a"},
		}
		body, _ := json.Marshal(host)
		c.Request, _ = http.NewRequest("", "", bytes.NewReader(body))

		svc := &mocks.HostService{}
		svc.On("SetHost", env, host).Return(nil)

		NewHostController(svc).CreateHost(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}
