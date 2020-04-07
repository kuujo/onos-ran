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

package codegen

// StorageType is the type of a storage backend
type StorageType string

// Config is the store generator configuration
type Config struct {
	Path    string        `yaml:"path,omitempty"`
	Package string        `yaml:"package,omitempty"`
	Stores  []StoreConfig `yaml:"stores"`
}

// StoreConfig is a store generator store configuration
type StoreConfig struct {
	Name    string              `yaml:"name,omitempty"`
	Package string              `yaml:"package,omitempty"`
	Message *MessageStoreConfig `yaml:"message,omitempty"`
	Storage StorageConfig       `yaml:"storage,omitempty"`
	Target  TargetConfig        `yaml:"target,omitempty"`
}

// MessageStoreConfig is a store generator configuration for message stores
type MessageStoreConfig struct {
	Type MessageTypeConfig `yaml:"type,omitempty"`
}

// MessageTypeConfig is a store generator data type configuration
type MessageTypeConfig struct {
	Package string `yaml:"package,omitempty"`
	Name    string `yaml:"name,omitempty"`
}

// StorageConfig is a store generator storage configuration
type StorageConfig struct {
	Type   StorageType `yaml:"type,omitempty"`
	Name   string      `yaml:"name,omitempty"`
	Cached bool        `yaml:"cached,omitempty"`
}

// TargetConfig is a store generator target configuration
type TargetConfig struct {
	Package string `yaml:"package,omitempty"`
}
