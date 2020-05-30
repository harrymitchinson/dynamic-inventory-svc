package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInventoryFile(t *testing.T) {

	t.Run("NoHosts", func(t *testing.T) {
		h := make([]Host, 0)

		f := NewInventoryFile(h)

		assert.Empty(t, f.All.Hosts)
		assert.Empty(t, f.All.Children)
	})

	t.Run("NoRoles", func(t *testing.T) {
		h := []Host{
			{
				Name:     "a",
				Hostname: "a",
				IP:       "10.0.0.1",
				Roles:    []string{},
			},
		}

		f := NewInventoryFile(h)

		assert.NotEmpty(t, f.All.Hosts)
		assert.Empty(t, f.All.Children)
	})

	t.Run("DifferentRoles", func(t *testing.T) {
		h := []Host{
			{
				Name:     "a",
				Hostname: "a",
				IP:       "10.0.0.1",
				Roles:    []string{"a"},
			},
			{
				Name:     "b",
				Hostname: "b",
				IP:       "10.0.0.2",
				Roles:    []string{"b"},
			},
		}

		f := NewInventoryFile(h)

		assert.Empty(t, f.All.Hosts)
		assert.Len(t, f.All.Children, 2)
	})

	t.Run("SameRoles", func(t *testing.T) {
		h := []Host{
			{
				Name:     "a",
				Hostname: "a",
				IP:       "10.0.0.1",
				Roles:    []string{"a"},
			},
			{
				Name:     "b",
				Hostname: "b",
				IP:       "10.0.0.2",
				Roles:    []string{"a"},
			},
		}

		f := NewInventoryFile(h)

		assert.Empty(t, f.All.Hosts)
		assert.Len(t, f.All.Children, 1)
	})

	t.Run("DuplicateHost", func(t *testing.T) {
		h := []Host{
			{
				Name:     "a",
				Hostname: "a",
				IP:       "10.0.0.1",
				Roles:    []string{"a"},
			},
			{
				Name:     "a",
				Hostname: "a",
				IP:       "10.0.0.1",
				Roles:    []string{"a"},
			},
		}

		f := NewInventoryFile(h)

		assert.Empty(t, f.All.Hosts)
		assert.Len(t, f.All.Children[h[0].Roles[0]].Hosts, 1)
	})

}
