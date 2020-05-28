# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/nb/a1/types/types.proto](#api/nb/a1/types/types.proto)
    - [CellID](#types.CellID)
    - [GroupID](#types.GroupID)
    - [PolicyID](#types.PolicyID)
    - [QosID](#types.QosID)
    - [ScopeIdentifier](#types.ScopeIdentifier)
    - [SliceID](#types.SliceID)
    - [UeID](#types.UeID)
  
    - [EnforcementReasonType](#types.EnforcementReasonType)
    - [EnforcementStatusType](#types.EnforcementStatusType)
    - [OperationStatus](#types.OperationStatus)
    - [PolicyQueryType](#types.PolicyQueryType)
    - [PolicyType](#types.PolicyType)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/nb/a1/types/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/nb/a1/types/types.proto



<a name="types.CellID"></a>

### CellID
CellID an identifier for a single cell

TODO it should be defined properly
The ECGI uniquely identifies the cell, while the C-RNTI uniquely identifies a UE within a gNB.
The ECGI&#43;C-RNTI pair constitutes a globally unique UE identifier






<a name="types.GroupID"></a>

### GroupID
GroupID an identifier for a group of UEs

TODO should be defined properly






<a name="types.PolicyID"></a>

### PolicyID
PolicyID Identifier of an A1 policy that is used in policy operations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |






<a name="types.QosID"></a>

### QosID
QoSID an identifier for QoS

TODO should be defined properly






<a name="types.ScopeIdentifier"></a>

### ScopeIdentifier
ScopeIdentifier Identifier of what the statements in the policy
applies to (UE, group of UEs, slice, QoS flow, network resource or combinations thereof).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ue_id | [UeID](#types.UeID) |  |  |
| group_id | [GroupID](#types.GroupID) |  |  |
| slice_id | [SliceID](#types.SliceID) |  |  |
| qos_id | [QosID](#types.QosID) |  |  |
| cell_id | [CellID](#types.CellID) |  |  |






<a name="types.SliceID"></a>

### SliceID
SliceID an identifier for a slice

TODO should be defined properly






<a name="types.UeID"></a>

### UeID
UeID an identifier for a single UE

TODO should be defined properly
The ECGI uniquely identifies the cell, while the C-RNTI uniquely identifies a UE within a gNB.
The ECGI&#43;C-RNTI pair constitutes a globally unique UE identifier





 


<a name="types.EnforcementReasonType"></a>

### EnforcementReasonType
EnforcementReasonType  represents the reason why notification is sent (e.g. why enforcement status has changed).

| Name | Number | Description |
| ---- | ------ | ----------- |
| SCOPE_NOT_APPLICABLE | 0 |  |
| STATEMENT_NOT_APPLICABLE | 1 |  |
| OTHER_REASON | 2 |  |



<a name="types.EnforcementStatusType"></a>

### EnforcementStatusType
EnforcementStatusType represents if a policy is enforced or not.

| Name | Number | Description |
| ---- | ------ | ----------- |
| ENFORCED | 0 |  |
| NOT_ENFORCED | 1 |  |
| UNDEFINED | 2 |  |



<a name="types.OperationStatus"></a>

### OperationStatus
OperationStatus status of performing an A1 methods

| Name | Number | Description |
| ---- | ------ | ----------- |
| SUCCESS | 0 |  |
| FAILED | 1 |  |



<a name="types.PolicyQueryType"></a>

### PolicyQueryType
PolicyQueryType type of a policy query (Query single policy, Query all policies, Query policy status)

| Name | Number | Description |
| ---- | ------ | ----------- |
| SINGLE_POLICY | 0 | get a single policy based on a given PolicyID |
| ALL_POLICIES | 1 | get all policies identities |
| POLICY_STATUS | 2 | get the policy status based on a given policyID |



<a name="types.PolicyType"></a>

### PolicyType
PolicyType

| Name | Number | Description |
| ---- | ------ | ----------- |
| QOS | 0 |  |
| TSP | 1 |  |


 

 

 



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

