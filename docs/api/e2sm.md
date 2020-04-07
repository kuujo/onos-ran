# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/sb/e2sm/e2sm.proto](#api/sb/e2sm/e2sm.proto)
    - [RicControlHeader](#interface.e2sm.RicControlHeader)
    - [RicControlMessage](#interface.e2sm.RicControlMessage)
    - [RicControlOutcome](#interface.e2sm.RicControlOutcome)
    - [RicEventTriggerDefinition](#interface.e2sm.RicEventTriggerDefinition)
    - [RicIndicationHeader](#interface.e2sm.RicIndicationHeader)
    - [RicIndicationMessage](#interface.e2sm.RicIndicationMessage)
  
    - [RanFunctionID](#interface.e2sm.RanFunctionID)
  
  
  

- [Scalar Value Types](#scalar-value-types)



<a name="api/sb/e2sm/e2sm.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/sb/e2sm/e2sm.proto
Copyright 2020-present Open Networking Foundation.

Licensed under the Apache License, Version 2.0 (the &#34;License&#34;);
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an &#34;AS IS&#34; BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


<a name="interface.e2sm.RicControlHeader"></a>

### RicControlHeader
RICcontrolHeader ::= OCTET STRING


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| messageType | [interface.e2.MessageType](#interface.e2.MessageType) |  |  |






<a name="interface.e2sm.RicControlMessage"></a>

### RicControlMessage
-- **************************************************************
-- Following IE defined in E2SM
-- **************************************************************
RICcontrolMessage ::= OCTET STRING


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rRMConfig | [interface.e2.RRMConfig](#interface.e2.RRMConfig) |  |  |
| hORequest | [interface.e2.HORequest](#interface.e2.HORequest) |  |  |
| cellConfigRequest | [interface.e2.CellConfigRequest](#interface.e2.CellConfigRequest) |  |  |






<a name="interface.e2sm.RicControlOutcome"></a>

### RicControlOutcome



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hOComplete | [interface.e2.HOComplete](#interface.e2.HOComplete) |  |  |
| hOFailure | [interface.e2.HOFailure](#interface.e2.HOFailure) |  |  |
| rRMConfigStatus | [interface.e2.RRMConfigStatus](#interface.e2.RRMConfigStatus) |  |  |






<a name="interface.e2sm.RicEventTriggerDefinition"></a>

### RicEventTriggerDefinition
-- **************************************************************
-- Following IE defined in E2SM
-- **************************************************************
RICeventTriggerDefinition ::= OCTET STRING


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| l2MeasConfig | [interface.e2.L2MeasConfig](#interface.e2.L2MeasConfig) |  |  |






<a name="interface.e2sm.RicIndicationHeader"></a>

### RicIndicationHeader
-- **************************************************************
-- Following IE defined in E2SM
-- **************************************************************
RICindicationHeader ::= OCTET STRING
RICindicationMessage ::= OCTET STRING


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| messageType | [interface.e2.MessageType](#interface.e2.MessageType) |  |  |






<a name="interface.e2sm.RicIndicationMessage"></a>

### RicIndicationMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| radioMeasReportPerUE | [interface.e2.RadioMeasReportPerUE](#interface.e2.RadioMeasReportPerUE) |  |  |
| uEAdmissionRequest | [interface.e2.UEAdmissionRequest](#interface.e2.UEAdmissionRequest) |  |  |
| uEReleaseInd | [interface.e2.UEReleaseInd](#interface.e2.UEReleaseInd) |  |  |





 


<a name="interface.e2sm.RanFunctionID"></a>

### RanFunctionID
-- **************************************************************
-- Following IE defined in E2SM
-- **************************************************************
RANfunctionDefinition ::= OCTET STRING

RANfunctionID ::= INTEGER (0..4095)

RANfunctionRevision ::= INTEGER (0..4095)

| Name | Number | Description |
| ---- | ------ | ----------- |
| RAN_FUNCTION_ID_INVALID | 0 |  |
| RAN_FUNCTION_ID_L2_MEAS | 1 |  |


 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

