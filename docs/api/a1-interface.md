# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/nb/a1/a1.proto](#api/nb/a1/a1.proto)
    - [AllPolicyResponse](#a1.AllPolicyResponse)
    - [CreateOrUpdateRequest](#a1.CreateOrUpdateRequest)
    - [CreateOrUpdateResponse](#a1.CreateOrUpdateResponse)
    - [DeleteRequest](#a1.DeleteRequest)
    - [DeleteResponse](#a1.DeleteResponse)
    - [NotifyRequest](#a1.NotifyRequest)
    - [NotifyResponse](#a1.NotifyResponse)
    - [Policy](#a1.Policy)
    - [PolicyStatement](#a1.PolicyStatement)
    - [PolicyStatus](#a1.PolicyStatus)
    - [PolicyStatusResponse](#a1.PolicyStatusResponse)
    - [ProblemDetails](#a1.ProblemDetails)
    - [QueryRequest](#a1.QueryRequest)
    - [QueryResponse](#a1.QueryResponse)
    - [SinglePolicyResponse](#a1.SinglePolicyResponse)
  
    - [a1](#a1.a1)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/nb/a1/a1.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/nb/a1/a1.proto



<a name="a1.AllPolicyResponse"></a>

### AllPolicyResponse
AllPolicyResponse representation of a A1 All policy response which contains list of all policy IDs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_id | [types.PolicyID](#types.PolicyID) | repeated |  |






<a name="a1.CreateOrUpdateRequest"></a>

### CreateOrUpdateRequest
CreateOrUpdateRequest a request to create or update a policy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_id | [types.PolicyID](#types.PolicyID) |  |  |
| policy | [Policy](#a1.Policy) |  |  |






<a name="a1.CreateOrUpdateResponse"></a>

### CreateOrUpdateResponse
CreateOrUpdateResponse response to a CreateOrUpdateRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy | [Policy](#a1.Policy) |  |  |
| status | [types.OperationStatus](#types.OperationStatus) |  |  |
| problem_details | [ProblemDetails](#a1.ProblemDetails) |  |  |






<a name="a1.DeleteRequest"></a>

### DeleteRequest
DeleteRequest a request to delete a policy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_id | [types.PolicyID](#types.PolicyID) |  |  |






<a name="a1.DeleteResponse"></a>

### DeleteResponse
DeleteResponse a response to a DeleteRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [types.OperationStatus](#types.OperationStatus) |  |  |
| problem_details | [ProblemDetails](#a1.ProblemDetails) |  |  |






<a name="a1.NotifyRequest"></a>

### NotifyRequest
NotifyRequest a request to get updates about changes in the policy enforcement status for an A1 policy;






<a name="a1.NotifyResponse"></a>

### NotifyResponse
NotifyResponse a notification response contains the information about changes and causes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_status | [PolicyStatus](#a1.PolicyStatus) |  |  |






<a name="a1.Policy"></a>

### Policy
Policy an A1 policy which contains a scope identifier and one or more policy statements


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| scope_id | [types.ScopeIdentifier](#types.ScopeIdentifier) |  |  |
| policy_statement | [PolicyStatement](#a1.PolicyStatement) | repeated |  |






<a name="a1.PolicyStatement"></a>

### PolicyStatement
PolicyStatement 	Expression of a directive in an A1 policy
that is related to policy objectives and/or policy resources
and are to be applied to/for the entities identified by the scope identifier.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [types.PolicyType](#types.PolicyType) |  |  |
| qos_objectives | [qos.QosObjectives](#a1.qos.QosObjectives) |  | QosObjectives |
| tsp_resources | [tsp.TspResources](#a1.tsp.TspResources) |  | TspResources Traffic steering optimization |






<a name="a1.PolicyStatus"></a>

### PolicyStatus
PolicyStatus representation of a A1  policy enforcement status and reasons


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| enforce_status | [types.EnforcementStatusType](#types.EnforcementStatusType) |  |  |
| enforce_reason | [types.EnforcementReasonType](#types.EnforcementReasonType) |  |  |






<a name="a1.PolicyStatusResponse"></a>

### PolicyStatusResponse
PolicyStatusResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy_status | [PolicyStatus](#a1.PolicyStatus) |  |  |






<a name="a1.ProblemDetails"></a>

### ProblemDetails
ProblemDetails


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) |  |  |
| detail | [string](#string) |  | TODO some other fields should be added |






<a name="a1.QueryRequest"></a>

### QueryRequest
QueryRequest a request to query about one or more than one policy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [types.PolicyQueryType](#types.PolicyQueryType) |  | Policy Type |
| policy_id | [types.PolicyID](#types.PolicyID) |  | PolicyID |






<a name="a1.QueryResponse"></a>

### QueryResponse
QueryResponse a response to a query request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| single_policy_response | [SinglePolicyResponse](#a1.SinglePolicyResponse) |  |  |
| all_policy_response | [AllPolicyResponse](#a1.AllPolicyResponse) |  |  |
| policy_status_response | [PolicyStatusResponse](#a1.PolicyStatusResponse) |  |  |
| status | [types.OperationStatus](#types.OperationStatus) |  |  |
| problem_details | [ProblemDetails](#a1.ProblemDetails) |  |  |






<a name="a1.SinglePolicyResponse"></a>

### SinglePolicyResponse
SinglePolicyResponse representation of a A1 single policy response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| policy | [Policy](#a1.Policy) |  |  |





 

 

 


<a name="a1.a1"></a>

### a1
a1 A1 policy based service

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateOrUpdate | [CreateOrUpdateRequest](#a1.CreateOrUpdateRequest) | [CreateOrUpdateResponse](#a1.CreateOrUpdateResponse) | CreateOrUpdate creates an A1 policy |
| Query | [QueryRequest](#a1.QueryRequest) | [QueryResponse](#a1.QueryResponse) stream | Query queries about one or more than one A1 policies |
| Delete | [DeleteRequest](#a1.DeleteRequest) | [DeleteResponse](#a1.DeleteResponse) | Delete deletes an A1 policy |
| Notify | [NotifyRequest](#a1.NotifyRequest) stream | [NotifyResponse](#a1.NotifyResponse) stream | Notify notify about an enforcement status change of a policy between &#39;enforced&#39; and &#39;not enforced&#39;. |

 



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

