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

package orane2

//#cgo CFLAGS: -I. -D_DEFAULT_SOURCE -DASN_DISABLE_OER_SUPPORT
//#cgo LDFLAGS: -lm
//#include <stdio.h>
//#include <stdlib.h>
//#include <assert.h>
//#include "RICindication.h"
//#include "Criticality.h"
//#include "ProtocolIE-Field.h"
//extern int consumeBytesCb(void* p0, uint32_t p1, void* p2);
import "C"
import (
	"fmt"
	"unsafe"
)

// This just dumps the XML to stdout
func dumpXer(ri *C.RICindication_t) error {
	ret, err := C.xer_fprint(nil, &C.asn_DEF_RICindication, unsafe.Pointer(ri))
	if err != nil {
		return err
	}
	if ret == -1 {
		return fmt.Errorf("error encoding. %d", int(ret))
	}
	return nil
}

// EncodeRicIndXer Encodes the XML to a byte[]
func EncodeRicIndXer(key int, ri *C.RICindication_t) ([]byte, error) {
	keyCint := C.int(key)

	// bytes get pushed back through this channel
	go func() {
		for d := range cbChan {
			chanMutex.Lock()
			if bytes, ok := responseMap[d.key]; ok {
				responseMap[d.key] = append(bytes, d.data...)
			} else {
				responseMap[d.key] = d.data
			}
			chanMutex.Unlock()
		}
	}()

	xerCbF := (*C.asn_app_consume_bytes_f)(C.consumeBytesCb)
	encRetVal, err := C.xer_encode(&C.asn_DEF_RICindication, unsafe.Pointer(ri),
		C.XER_F_BASIC, xerCbF, unsafe.Pointer(&keyCint))
	if err != nil {
		return nil, err
	}
	if encRetVal.encoded == -1 {
		fmt.Printf("error on %v\n", *encRetVal.failed_type)

		return nil, fmt.Errorf("error encoding. Name: %v Tag: %v",
			C.GoString(encRetVal.failed_type.name),
			C.GoString(encRetVal.failed_type.name))
	}
	chanMutex.RLock()
	defer chanMutex.RUnlock()
	return responseMap[key], err
}
