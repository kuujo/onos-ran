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
import "C"
import (
	"encoding/binary"
	"unsafe"
)

// BuildRicIndication --
func buildRicIndication() (*C.RICindication_t, error) {
	ri := C.RICindication_t{
		protocolIEs: C.ProtocolIE_Container_1544P6_t{},
	}

	vpr, v := buildRICrequestID(1, 2)
	var riIe *C.RICindication_IEs_t
	var err error
	if riIe, err = buildRicIndicationIEValue(
		C.ProtocolIE_ID_id_RICrequestID, C.Criticality_reject, vpr, v); err != nil {
		return nil, err
	}
	if _, err = C.asn_sequence_add(unsafe.Pointer(&ri.protocolIEs), unsafe.Pointer(riIe)); err != nil {
		return nil, err
	}

	vpr2, v2 := buildRANfunctionID(C.RANfunctionID_t(1))
	if riIe, err = buildRicIndicationIEValue(
		C.ProtocolIE_ID_id_RANfunctionID, C.Criticality_reject, vpr2, v2); err != nil {
		return nil, err
	}
	if _, err = C.asn_sequence_add(unsafe.Pointer(&ri.protocolIEs), unsafe.Pointer(riIe)); err != nil {
		return nil, err
	}

	return &ri, nil
}

func buildRicIndicationIEValue(id C.ProtocolIE_ID_t,
	crit C.Criticality_t, vpr C.RICindication_IEs__value_PR,
	value C.union_RICcontrolRequest_IEs__value_u) (*C.RICindication_IEs_t, error) {

	riIe := C.RICindication_IEs_t{
		id:          id,
		criticality: crit,
		value: C.struct_RICindication_IEs__value{
			present: vpr,
			choice:  value,
		},
	}

	return &riIe, nil
}

func buildRICrequestID(requestorID int, instanceID int) (C.RICindication_IEs__value_PR, C.union_RICcontrolRequest_IEs__value_u) {
	reqID := C.RICrequestID_t{ricRequestorID: C.long(requestorID), ricInstanceID: C.long(instanceID)}
	var choice C.union_RICcontrolRequest_IEs__value_u
	binary.LittleEndian.PutUint64(choice[0:], uint64(reqID.ricRequestorID))
	binary.LittleEndian.PutUint64(choice[8:], uint64(reqID.ricInstanceID))

	return C.RICindication_IEs__value_PR_RICrequestID, choice
}

func buildRANfunctionID(funcID C.RANfunctionID_t) (C.RICindication_IEs__value_PR, C.union_RICcontrolRequest_IEs__value_u) {
	var choice C.union_RICcontrolRequest_IEs__value_u
	binary.LittleEndian.PutUint64(choice[0:], uint64(funcID))

	return C.RICindication_IEs__value_PR_RANfunctionID, choice
}
