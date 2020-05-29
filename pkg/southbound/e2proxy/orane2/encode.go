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
	"sync"
	"unsafe"
)

type cbData struct {
	data []byte
	key  int
}

var cbChan = make(chan cbData)
var responseMap = make(map[int][]byte)
var chanMutex = sync.RWMutex{}

// Callback function of type C.asn_app_consume_bytes_f
//export consumeBytesCb
func consumeBytesCb(buf unsafe.Pointer, size C.uint32_t, key unsafe.Pointer) C.int {
	bytes := make([]byte, int(size))
	for i := 0; i < int(size); i++ {
		addr := *(*C.uchar)(unsafe.Pointer(uintptr(buf) + uintptr(i)))
		bytes = append(bytes, byte(addr))
	}
	cbChan <- cbData{
		data: bytes,
		key:  int(*(*C.int)(key)),
	}
	return C.int(size)
}
