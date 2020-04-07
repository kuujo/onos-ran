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
	"path"
)

func generateStores(options Options) error {
	for _, store := range options.Stores {
		if store.Impl.Storage.Cached {
			if err := generateCachedCRUDInterface(store); err != nil {
				return err
			}
			if err := generateCachedCRUDImpl(store); err != nil {
				return err
			}
		} else {
			return errors.New("unsupported storage configuration")
		}
	}
	return nil
}

func generateCachedCRUDInterface(options StoreOptions) error {
	return generateTemplate(getTemplate("message_interface.tpl"), path.Join(options.Interface.Location.Path, options.Interface.Location.File), options)
}

func generateCachedCRUDImpl(options StoreOptions) error {
	return generateTemplate(getTemplate("message_store.tpl"), path.Join(options.Impl.Location.Path, options.Impl.Location.File), options)
}
