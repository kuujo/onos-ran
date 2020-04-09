# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/sb/e2ap/e2ap.proto](#api/sb/e2ap/e2ap.proto)
    - [RicControlRequest](#interface.e2ap.RicControlRequest)
    - [RicControlResponse](#interface.e2ap.RicControlResponse)
    - [RicIndication](#interface.e2ap.RicIndication)
    - [RicSubscriptionDetails](#interface.e2ap.RicSubscriptionDetails)
    - [RicSubscriptionRequest](#interface.e2ap.RicSubscriptionRequest)
  
  
  
    - [E2AP](#interface.e2ap.E2AP)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api/sb/e2ap/e2ap.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/sb/e2ap/e2ap.proto
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


<a name="interface.e2ap.RicControlRequest"></a>

### RicControlRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hdr | [interface.e2sm.RicControlHeader](#interface.e2sm.RicControlHeader) |  |  |
| msg | [interface.e2sm.RicControlMessage](#interface.e2sm.RicControlMessage) |  |  |






<a name="interface.e2ap.RicControlResponse"></a>

### RicControlResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hdr | [interface.e2sm.RicControlHeader](#interface.e2sm.RicControlHeader) |  |  |
| msg | [interface.e2sm.RicControlOutcome](#interface.e2sm.RicControlOutcome) |  |  |






<a name="interface.e2ap.RicIndication"></a>

### RicIndication



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hdr | [interface.e2sm.RicIndicationHeader](#interface.e2sm.RicIndicationHeader) |  |  |
| msg | [interface.e2sm.RicIndicationMessage](#interface.e2sm.RicIndicationMessage) |  |  |






<a name="interface.e2ap.RicSubscriptionDetails"></a>

### RicSubscriptionDetails



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ricEventTriggerDefinition | [interface.e2sm.RicEventTriggerDefinition](#interface.e2sm.RicEventTriggerDefinition) |  |  |






<a name="interface.e2ap.RicSubscriptionRequest"></a>

### RicSubscriptionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ranFunctionID | [interface.e2sm.RanFunctionID](#interface.e2sm.RanFunctionID) |  |  |
| ricSubscriptionDetails | [RicSubscriptionDetails](#interface.e2ap.RicSubscriptionDetails) |  |  |





 

 

 


<a name="interface.e2ap.E2AP"></a>

### E2AP


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RicSubscription | [RicSubscriptionRequest](#interface.e2ap.RicSubscriptionRequest) stream | [RicIndication](#interface.e2ap.RicIndication) stream |  |
| RicControl | [RicControlRequest](#interface.e2ap.RicControlRequest) stream | [RicControlResponse](#interface.e2ap.RicControlResponse) stream |  |

 



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

