// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	configlib "github.com/onosproject/onos-lib-go/pkg/config"
)

const defaultPartitions = 16

var config *Config

// Config is the onos-config configuration
type Config struct {
	// Mastership is the mastership configuration
	Mastership MastershipConfig `yaml:"mastership,omitempty"`
	// Atomix is the Atomix configuration
	Atomix atomix.Config `yaml:"atomix,omitempty"`
}

// MastershipConfig is the configuration for store mastership
type MastershipConfig struct {
	// Partitions is the number of store partitions
	Partitions int `yaml:"partitions,omitempty"`
}

// GetPartitions returns the number of store partitions
func (c MastershipConfig) GetPartitions() int {
	partitions := c.Partitions
	if partitions == 0 {
		partitions = defaultPartitions
	}
	return partitions
}

// GetConfig gets the onos-config configuration
func GetConfig() (Config, error) {
	if config == nil {
		config = &Config{}
		if err := configlib.Load(config); err != nil {
			return Config{}, err
		}
	}
	return *config, nil
}

// GetConfigOrDie gets the onos-config configuration or panics
func GetConfigOrDie() Config {
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}
	return config
}
