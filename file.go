package inventory

// File is the top level inventory file.
type File struct {
	All All `json:"all"`
}

// All is a section in the inventory file.
type All struct {
	Hosts    map[string]*string `json:"hosts"`
	Children map[string]Child   `json:"children"`
}

// Child is group of hosts in the inventory file.
type Child struct {
	Hosts map[string]*string `json:"hosts"`
}

// NewInventoryFile creates a File from a []Host
func NewInventoryFile(hosts []Host) *File {
	file := File{
		All: All{
			Hosts:    map[string]*string{},
			Children: map[string]Child{},
		},
	}

	for _, host := range hosts {
		if len(host.Roles) == 0 {
			file.All.Hosts[host.Hostname] = nil
			continue
		}
		for _, role := range host.Roles {
			if x, exists := file.All.Children[role]; exists {
				x.Hosts[host.Hostname] = nil
			} else {
				file.All.Children[role] = Child{
					Hosts: map[string]*string{
						host.Hostname: nil,
					},
				}
			}
		}
	}

	return &file
}
