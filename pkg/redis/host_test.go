package redis

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	inventory "github.com/harrymitchinson/dynamic-inventory-svc"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type helper struct {
	m *miniredis.Miniredis
	r *redis.Client
}

func (h helper) Close() {
	h.r.Close()
	h.m.Close()
}

func getTestHelper(t *testing.T) (helper, HostService) {
	l := zaptest.NewLogger(t, zaptest.Level(zap.ErrorLevel))

	m, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	r := redis.NewClient(&redis.Options{
		Addr: m.Addr(),
	})

	if err != nil {
		panic(err)
	}

	h := helper{
		m: m,
		r: r,
	}
	svc := NewHostService(r, l)
	return h, svc

}

func TestGetHosts(t *testing.T) {
	h, svc := getTestHelper(t)
	defer h.Close()

	environment := inventory.EnvironmentDev
	host := inventory.Host{
		Name:     "a",
		Hostname: "a",
		IP:       "10.0.0.1",
		Roles:    []string{"a"},
	}
	key := fmt.Sprintf(KeyFormatHostsFull, environment, host.Name)
	value, _ := host.MarshalBinary()
	h.m.Set(key, string(value))

	hosts, err := svc.GetHosts(inventory.EnvironmentDev)

	if assert.NoError(t, err) {
		assert.NotEmpty(t, hosts)
	}
}

func TestSetHost(t *testing.T) {
	h, svc := getTestHelper(t)
	defer h.Close()

	environment := inventory.EnvironmentDev
	host := &inventory.Host{
		Name:     "a",
		Hostname: "a",
		IP:       "10.0.0.1",
		Roles:    []string{"a"},
	}

	err := svc.SetHost(environment, host)

	if assert.NoError(t, err) {
		key := fmt.Sprintf(KeyFormatHostsFull, environment, host.Name)
		expected, _ := host.MarshalBinary()
		h.m.CheckGet(t, key, string(expected))
	}
}
