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
	// Deprecated: use StoreConfig instead
	Mastership MastershipStoreConfig `yaml:"mastership,omitempty"`
	// Atomix is the Atomix configuration
	Atomix atomix.Config `yaml:"atomix,omitempty"`
	// Stores is the store configurations
	Stores StoresConfig `yaml:"stores,omitempty"`
}

// StoresConfig is the configuration for stores
type StoresConfig struct {
	Mastership MastershipStoreConfig `yaml:"mastership,omitempty"`
	Requests   RequestsStoreConfig   `yaml:"requests,omitempty"`
}

// MastershipStoreConfig is the configuration for store mastership
type MastershipStoreConfig struct {
	// Partitions is the number of store partitions
	Partitions int `yaml:"partitions,omitempty"`
}

// WithConfig passes in a predefined Config to use for testing
func WithConfig(newConfig *Config) {
	config = newConfig
}

// GetPartitions returns the number of store partitions
func (c MastershipStoreConfig) GetPartitions() int {
	partitions := c.Partitions
	if partitions == 0 {
		partitions = defaultPartitions
	}
	return partitions
}

// RequestsStoreConfig is the requests store configuration
type RequestsStoreConfig struct {
	// SyncBackups configures the number of synchronous backups
	SyncBackups *int `yaml:"syncBackups,omitempty"`
}

// GetSyncBackups returns the number of synchronous backups
func (c RequestsStoreConfig) GetSyncBackups() int {
	backups := c.SyncBackups
	if backups == nil {
		return 1
	}
	return *backups
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
