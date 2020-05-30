package inventory

import (
	"encoding/json"
)

// Host is a simple representation of a server, used for building the inventory file.
type Host struct {
	Name     string   `json:"name" validate:"required,hostname"`
	Hostname string   `json:"hostname" validate:"required,hostname"`
	IP       string   `json:"ip" validate:"required,ip"`
	Roles    []string `json:"roles,omitempty"`
}

// HostService is the interface to interact with Hosts.
type HostService interface {
	GetHosts(e Environment) ([]Host, error)
	SetHost(e Environment, h *Host) error
}

// MarshalBinary marshals the Host to a []byte.
func (h *Host) MarshalBinary() ([]byte, error) {
	return json.Marshal(h)
}

// UnmarshalBinary unmarshals the Host from a []byte.
func (h *Host) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &h)
}
