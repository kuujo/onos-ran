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

// EncRicIndToBuffer encodes a RICIndication to a []byte
func EncRicIndToBuffer(ri *C.RICindication_t) ([]byte, error) {
	perBuf := C.malloc(C.sizeof_uchar * 1024) // C allocated pointer
	defer C.free(perBuf)
	encRetVal, err := C.uper_encode_to_buffer(
		&C.asn_DEF_RICindication, nil, unsafe.Pointer(ri),
		perBuf, C.ulong(1024))
	if err != nil {
		return nil, err
	}
	if encRetVal.encoded == -1 {
		fmt.Printf("error on %v\n", *encRetVal.failed_type)

		return nil, fmt.Errorf("error encoding. Name: %v Tag: %v",
			C.GoString(encRetVal.failed_type.name),
			C.GoString(encRetVal.failed_type.xml_tag))
	}
	bytes := make([]byte, encRetVal.encoded)
	for i := 0; i < int(encRetVal.encoded); i++ {
		b := *(*C.uchar)(unsafe.Pointer(uintptr(perBuf) + uintptr(i)))
		bytes[i] = byte(b)
	}
	return bytes, nil
}
