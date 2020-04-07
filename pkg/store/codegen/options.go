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

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// Options is code generator options
type Options struct {
	Stores map[string]StoreOptions
}

// LocationOptions is the location of a code file
type LocationOptions struct {
	Path string
	File string
}

// PackageOptions is the package for a code file
type PackageOptions struct {
	Name  string
	Path  string
	Alias string
}

// StoreOptions is the options for a store
type StoreOptions struct {
	DataType  DataTypeOptions
	Interface InterfaceOptions
	Impl      ImplOptions
}

// InterfaceOptions is the options for a store interface
type InterfaceOptions struct {
	Location LocationOptions
	Package  PackageOptions
	Type     TypeOptions
}

// ImplOptions is the options for a store implementation
type ImplOptions struct {
	Location LocationOptions
	Package  PackageOptions
	Type     TypeOptions
	Storage  StorageOptions
}

// TypeOptions is the options for a store type
type TypeOptions struct {
	Name string
}

// DataTypeOptions is the options for a store data type
type DataTypeOptions struct {
	Package PackageOptions
	Name    string
}

// StorageOptions is the options for a store primitive
type StorageOptions struct {
	Type      string
	Primitive string
	Cached    bool
}

func getOptionsFromConfig(config Config) (Options, error) {
	options := Options{
		Stores: make(map[string]StoreOptions),
	}

	basePath, err := filepath.Abs(config.Path)
	if err != nil {
		return Options{}, err
	}

	for _, store := range config.Stores {
		_, ok := options.Stores[store.Name]
		if !ok {
			pkgPath := store.Target.Package
			if pkgPath == "" {
				pkgPath = fmt.Sprintf("%s/%s", config.Package, store.Name)
			}
			if pkgPath == "" {
				return Options{}, errors.New("no target package configured")
			}
			srcPath := path.Join(basePath, strings.Replace(pkgPath, "github.com/onosproject/onos-ric/", "", 1))

			if store.Message == nil {
				return Options{}, fmt.Errorf("no message type configured for %s", store.Name)
			}

			pkg := PackageOptions{
				Name:  path.Base(pkgPath),
				Path:  pkgPath,
				Alias: store.Name,
			}

			storeOptions := StoreOptions{
				DataType: DataTypeOptions{
					Package: PackageOptions{
						Name:  path.Base(store.Message.Type.Package),
						Path:  store.Message.Type.Package,
						Alias: path.Base(store.Message.Type.Package),
					},
					Name: store.Message.Type.Name,
				},
				Interface: InterfaceOptions{
					Location: LocationOptions{
						Path: srcPath,
						File: "interface.go",
					},
					Package: pkg,
					Type: TypeOptions{
						Name: "Store",
					},
				},
				Impl: ImplOptions{
					Location: LocationOptions{
						Path: srcPath,
						File: "store.go",
					},
					Package: pkg,
					Type: TypeOptions{
						Name: "store",
					},
					Storage: StorageOptions{
						Primitive: store.Storage.Name,
						Type:      string(store.Storage.Type),
						Cached:    store.Storage.Cached,
					},
				},
			}
			options.Stores[store.Name] = storeOptions
		}
	}
	return options, nil
}
