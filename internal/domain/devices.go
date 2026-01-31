package domain

type Config struct {
	Devices []Device `yaml:"devices"`
}

type Device struct {
	Name     string    `yaml:"name"`
	Profiles []Profile `yaml:"profiles"`
}

// AllProfiles returns a flat list of all profiles from all devices
func (c *Config) AllProfiles() []Profile {
	var result []Profile

	for _, device := range c.Devices {
		result = append(result, device.Profiles...)
	}

	return result
}

// AllGamers returns a flat list of all players from all profiles of all devices
func (c *Config) AllGamers() []*Gamer {
	var result []*Gamer

	for _, device := range c.Devices {
		for _, profile := range device.Profiles {
			for i := range profile.Gamer {
				result = append(result, &profile.Gamer[i])
			}
		}
	}

	return result
}
