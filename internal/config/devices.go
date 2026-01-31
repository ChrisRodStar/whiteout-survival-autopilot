package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/repository"
)

// LoadDeviceConfig reads the device configuration YAML file and deserializes it into the domain.Config structure.
func LoadDeviceConfig(devicesFile string, repo repository.StateRepository) (*domain.Config, error) {
	ctx := context.Background()

	// 📄 Load devices.yaml
	devicesData, err := os.ReadFile(filepath.Clean(devicesFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read devices.yaml: %w", err)
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(devicesData, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal devices.yaml: %w", err)
	}

	// 🧠 Load state from repository
	state, err := repo.LoadState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load state.yaml from repo: %w", err)
	}

	// 🔍 Index state by gamer.ID
	stateMap := make(map[int]domain.Gamer)
	for _, g := range state.Gamers {
		stateMap[g.ID] = g
	}

	// 🔁 Merge by ID and sort for stable order
	for dIdx := range cfg.Devices {
		for pIdx := range cfg.Devices[dIdx].Profiles {
			// 🔄 Merge state for each player
			for gIdx, gamer := range cfg.Devices[dIdx].Profiles[pIdx].Gamer {
				if full, ok := stateMap[gamer.ID]; ok {
					cfg.Devices[dIdx].Profiles[pIdx].Gamer[gIdx] = full
				}
			}

			// 🔡 Sort players by Nickname
			sort.Sort(domain.Gamers(cfg.Devices[dIdx].Profiles[pIdx].Gamer))
		}

		// 📧 Sort profiles by Email
		sort.Sort(domain.Profiles(cfg.Devices[dIdx].Profiles))
	}

	return &cfg, nil
}
