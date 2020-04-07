# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/nb/c1-interface.proto](#api/nb/c1-interface.proto)
    - [C1CandScell](#interface.c1.C1CandScell)
    - [C1CellConfigAttribute](#interface.c1.C1CellConfigAttribute)
    - [C1ECGI](#interface.c1.C1ECGI)
    - [C1HandoverRequest](#interface.c1.C1HandoverRequest)
    - [C1PciArfcn](#interface.c1.C1PciArfcn)
    - [C1RNIBCell](#interface.c1.C1RNIBCell)
    - [C1RNIBCells](#interface.c1.C1RNIBCells)
    - [C1RNIBLink](#interface.c1.C1RNIBLink)
    - [C1RNIBLinkID](#interface.c1.C1RNIBLinkID)
    - [C1RNIBLinks](#interface.c1.C1RNIBLinks)
    - [C1RNIBUE](#interface.c1.C1RNIBUE)
    - [C1RNIBUEs](#interface.c1.C1RNIBUEs)
    - [C1RRMConfigAttribute](#interface.c1.C1RRMConfigAttribute)
    - [C1RRMConfiguration](#interface.c1.C1RRMConfiguration)
    - [C1RadioMeasReportPerUeAttribute](#interface.c1.C1RadioMeasReportPerUeAttribute)
    - [C1RadioRepPerServCellAttribute](#interface.c1.C1RadioRepPerServCellAttribute)
    - [C1RequestMessage](#interface.c1.C1RequestMessage)
    - [C1RequestMessageHeader](#interface.c1.C1RequestMessageHeader)
    - [C1RequestMessagePayload](#interface.c1.C1RequestMessagePayload)
    - [C1ResponseMessage](#interface.c1.C1ResponseMessage)
    - [C1ResponseMessageHeader](#interface.c1.C1ResponseMessageHeader)
    - [C1ResponseMessagePayload](#interface.c1.C1ResponseMessagePayload)
    - [ChannelQuality](#interface.c1.ChannelQuality)
    - [ECGI](#interface.c1.ECGI)
    - [HandOverRequest](#interface.c1.HandOverRequest)
    - [HandOverResponse](#interface.c1.HandOverResponse)
    - [RadioPowerRequest](#interface.c1.RadioPowerRequest)
    - [RadioPowerResponse](#interface.c1.RadioPowerResponse)
    - [StationInfo](#interface.c1.StationInfo)
    - [StationLinkInfo](#interface.c1.StationLinkInfo)
    - [StationLinkListRequest](#interface.c1.StationLinkListRequest)
    - [StationListRequest](#interface.c1.StationListRequest)
    - [UELinkInfo](#interface.c1.UELinkInfo)
    - [UELinkListRequest](#interface.c1.UELinkListRequest)
  
    - [C1MessageType](#interface.c1.C1MessageType)
    - [C1RNIBType](#interface.c1.C1RNIBType)
    - [C1XICICPA](#interface.c1.C1XICICPA)
    - [StationPowerOffset](#interface.c1.StationPowerOffset)
  
  
    - [C1InterfaceService](#interface.c1.C1InterfaceService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api/nb/c1-interface.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/nb/c1-interface.proto



<a name="interface.c1.C1CandScell"></a>

### C1CandScell



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pci | [string](#string) |  |  |
| earfcnDl | [string](#string) |  |  |






<a name="interface.c1.C1CellConfigAttribute"></a>

### C1CellConfigAttribute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [C1ECGI](#interface.c1.C1ECGI) |  |  |
| pci | [string](#string) |  |  |
| candScells | [C1CandScell](#interface.c1.C1CandScell) | repeated |  |
| earfcnDl | [string](#string) |  |  |
| earfcnUl | [string](#string) |  |  |
| rbsPerTtiDl | [string](#string) |  |  |
| rbsPerTtiUl | [string](#string) |  |  |
| numTxAntenna | [string](#string) |  |  |
| duplexMode | [string](#string) |  |  |
| maxNumConnectedUes | [string](#string) |  |  |
| maxNumConnectedBearers | [string](#string) |  |  |
| maxNumUesSchedPerTtiDl | [string](#string) |  |  |
| maxNumUesSchedPerTtiUl | [string](#string) |  |  |
| dlfsSchedEnable | [string](#string) |  |  |






<a name="interface.c1.C1ECGI"></a>

### C1ECGI



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| plmnId | [string](#string) |  |  |
| ecid | [string](#string) |  |  |






<a name="interface.c1.C1HandoverRequest"></a>

### C1HandoverRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| links | [C1RNIBLinks](#interface.c1.C1RNIBLinks) |  | UE in links[index] should be moved from srcCells[index] to dstCell[index] |
| srcCells | [C1RNIBCells](#interface.c1.C1RNIBCells) |  |  |
| dstCells | [C1RNIBCells](#interface.c1.C1RNIBCells) |  |  |






<a name="interface.c1.C1PciArfcn"></a>

### C1PciArfcn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pci | [string](#string) |  |  |
| earfcnDl | [string](#string) |  |  |






<a name="interface.c1.C1RNIBCell"></a>

### C1RNIBCell



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [C1ECGI](#interface.c1.C1ECGI) |  | ID |
| cellConfiguration | [C1CellConfigAttribute](#interface.c1.C1CellConfigAttribute) |  | attributes of R-NIB Cell for Cell Configuration |
| rrmConfiguration | [C1RRMConfigAttribute](#interface.c1.C1RRMConfigAttribute) |  | for RRM configuration |






<a name="interface.c1.C1RNIBCells"></a>

### C1RNIBCells



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rNIBCells | [C1RNIBCell](#interface.c1.C1RNIBCell) | repeated |  |






<a name="interface.c1.C1RNIBLink"></a>

### C1RNIBLink



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| linkId | [C1RNIBLinkID](#interface.c1.C1RNIBLinkID) |  | ID |
| radioMeasReportPerUe | [C1RadioMeasReportPerUeAttribute](#interface.c1.C1RadioMeasReportPerUeAttribute) |  | attributes of R-NIB Link |






<a name="interface.c1.C1RNIBLinkID"></a>

### C1RNIBLinkID



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [C1ECGI](#interface.c1.C1ECGI) |  |  |
| imsi | [string](#string) |  |  |






<a name="interface.c1.C1RNIBLinks"></a>

### C1RNIBLinks



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rNIBLinks | [C1RNIBLink](#interface.c1.C1RNIBLink) | repeated |  |






<a name="interface.c1.C1RNIBUE"></a>

### C1RNIBUE



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| imsi | [string](#string) |  | ID |
| crnti | [string](#string) |  | attributes of R-NIB UE |
| sCell | [C1ECGI](#interface.c1.C1ECGI) |  |  |
| mmeUeS1apId | [string](#string) |  |  |
| enbUeS1apId | [string](#string) |  |  |






<a name="interface.c1.C1RNIBUEs"></a>

### C1RNIBUEs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rNIBUEs | [C1RNIBUE](#interface.c1.C1RNIBUE) | repeated |  |






<a name="interface.c1.C1RRMConfigAttribute"></a>

### C1RRMConfigAttribute
This attributed defined all elements as a list (repeated), but this list only has a single element for MWC demo -&gt; will be extended


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [C1ECGI](#interface.c1.C1ECGI) |  |  |
| crnti | [string](#string) | repeated | going to be FFFF or null, which means all crnti |
| pciArfcn | [C1PciArfcn](#interface.c1.C1PciArfcn) |  |  |
| pa | [C1XICICPA](#interface.c1.C1XICICPA) | repeated |  |
| startPrbDl | [string](#string) | repeated |  |
| endPrbDl | [string](#string) | repeated |  |
| subFrameBitmaskDl | [string](#string) | repeated |  |
| p0UePusch | [string](#string) | repeated |  |
| startPrbUl | [string](#string) | repeated |  |
| endPrbUl | [string](#string) | repeated |  |
| subFrameBitmaskUl | [string](#string) | repeated |  |






<a name="interface.c1.C1RRMConfiguration"></a>

### C1RRMConfiguration



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| targetCells | [C1RNIBCells](#interface.c1.C1RNIBCells) |  | targetCells[index] should have the offset pa which MLB app assigns; pa value in targetCell[index] is the new value to update on RNIBCell |






<a name="interface.c1.C1RadioMeasReportPerUeAttribute"></a>

### C1RadioMeasReportPerUeAttribute



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| radioRepPerServCell | [C1RadioRepPerServCellAttribute](#interface.c1.C1RadioRepPerServCellAttribute) | repeated |  |






<a name="interface.c1.C1RadioRepPerServCellAttribute"></a>

### C1RadioRepPerServCellAttribute
either serving cell or neighbor cell report


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [C1ECGI](#interface.c1.C1ECGI) |  |  |
| cqiHist | [string](#string) | repeated |  |
| riHist | [string](#string) | repeated |  |
| puschSinrHist | [string](#string) | repeated |  |
| pucchSinrHist | [string](#string) | repeated |  |






<a name="interface.c1.C1RequestMessage"></a>

### C1RequestMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [C1RequestMessageHeader](#interface.c1.C1RequestMessageHeader) |  |  |
| payload | [C1RequestMessagePayload](#interface.c1.C1RequestMessagePayload) |  |  |






<a name="interface.c1.C1RequestMessageHeader"></a>

### C1RequestMessageHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [C1MessageType](#interface.c1.C1MessageType) |  |  |






<a name="interface.c1.C1RequestMessagePayload"></a>

### C1RequestMessagePayload



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| requestedRNIBType | [C1RNIBType](#interface.c1.C1RNIBType) |  | for C1_GET_* |
| handoverRequest | [C1HandoverRequest](#interface.c1.C1HandoverRequest) |  | for C1_POST_HANDOVERS |
| rrmConfigurationRequest | [C1RRMConfiguration](#interface.c1.C1RRMConfiguration) |  | for C1_POST_RRMCONFIGURATION |






<a name="interface.c1.C1ResponseMessage"></a>

### C1ResponseMessage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [C1ResponseMessageHeader](#interface.c1.C1ResponseMessageHeader) |  |  |
| payload | [C1ResponseMessagePayload](#interface.c1.C1ResponseMessagePayload) |  |  |






<a name="interface.c1.C1ResponseMessageHeader"></a>

### C1ResponseMessageHeader



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [C1MessageType](#interface.c1.C1MessageType) |  |  |






<a name="interface.c1.C1ResponseMessagePayload"></a>

### C1ResponseMessagePayload



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| responseCode | [string](#string) |  | for POST_* |
| rNIBCells | [C1RNIBCells](#interface.c1.C1RNIBCells) |  | for GET_RNIBCELLS |
| rNIBUEs | [C1RNIBUEs](#interface.c1.C1RNIBUEs) |  | for GET_RNIBUES |
| rNIBLinks | [C1RNIBLinks](#interface.c1.C1RNIBLinks) |  | for GET_LINKS |






<a name="interface.c1.ChannelQuality"></a>

### ChannelQuality



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| targetEcgi | [ECGI](#interface.c1.ECGI) |  | Target stations&#39;s ID This target station can be either the serving station or the serving station&#39;s neighbor stations. |
| cqiHist | [uint32](#uint32) |  | CQI stands for Channel Quality Indicator in LTE, which ranges from 0 (out of range) to 15 (64 QAM and 948 Code rate) |






<a name="interface.c1.ECGI"></a>

### ECGI
station&#39;s unique ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| plmnid | [string](#string) |  | one of ECGI value |
| ecid | [string](#string) |  | one of ECGI value |






<a name="interface.c1.HandOverRequest"></a>

### HandOverRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  | UE&#39;s local ID in serving station |
| srcStation | [ECGI](#interface.c1.ECGI) |  | UE&#39;s source station ID - serving station |
| dstStation | [ECGI](#interface.c1.ECGI) |  | UE&#39;s destination station ID for handover - one of neighbor stations |






<a name="interface.c1.HandOverResponse"></a>

### HandOverResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  | empty - no response will be okay |






<a name="interface.c1.RadioPowerRequest"></a>

### RadioPowerRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.c1.ECGI) |  | target station&#39;s ID |
| offset | [StationPowerOffset](#interface.c1.StationPowerOffset) |  | target station&#39;s power offset to adjust transmission power |






<a name="interface.c1.RadioPowerResponse"></a>

### RadioPowerResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  | empty - no response will be okay |






<a name="interface.c1.StationInfo"></a>

### StationInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.c1.ECGI) |  | station&#39;s unique ID |
| maxNumConnectedUes | [uint32](#uint32) |  | station&#39;s maximum number of connected UEs - used for MLB |






<a name="interface.c1.StationLinkInfo"></a>

### StationLinkInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.c1.ECGI) |  | target station ID |
| neighborECGI | [ECGI](#interface.c1.ECGI) | repeated | list of neighbor stations&#39; ID |






<a name="interface.c1.StationLinkListRequest"></a>

### StationLinkListRequest
if ECGI is empty, stream all station links&#39; information


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.c1.ECGI) |  | ecgi - optional station identifier - list for all stations if not present |
| subscribe | [bool](#bool) |  | subscribe indicates whether to subscribe to events (e.g. ADD, UPDATE, and REMOVE) that occur after all stationlinks have been streamed to the client |






<a name="interface.c1.StationListRequest"></a>

### StationListRequest
if ecgi is empty, stream all stations&#39; information


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.c1.ECGI) |  | ecgi - optional station identifier - list for all stations if not present |
| subscribe | [bool](#bool) |  | subscribe indicates whether to subscribe to events (e.g. ADD, UPDATE, and REMOVE) that occur after all stations have been streamed to the client |






<a name="interface.c1.UELinkInfo"></a>

### UELinkInfo
It shows the link quality between the UE -- having crnti (1) and serviced by the station having ecgi (2) -- and the station -- having targetECGI (3).
Target station can be not only UE&#39;s serving station but also the station&#39;s neighbor stations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  | Both crnti and ecgi are used as a key in our store. UE&#39;s local ID in serving station |
| ecgi | [ECGI](#interface.c1.ECGI) |  | UE&#39;s serving station ID |
| channelQualities | [ChannelQuality](#interface.c1.ChannelQuality) | repeated | Channel quality values between the UE and its serving or neighbor stations. |
| imsi | [string](#string) |  | optional value: IMSI which is a global unique UE ID This value is from UEContextUpdate message The crnti is local ID only working in a serving base station. However, sometimes an app may need the unique ID Thus, add IMSI as an optional. |






<a name="interface.c1.UELinkListRequest"></a>

### UELinkListRequest
if crnti and ecgi are empty, stream all UE links&#39; information


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  | crnti - optional UE&#39;s local ID in serving station - list for all UE local IDs if not present |
| ecgi | [ECGI](#interface.c1.ECGI) |  | ecgi - optional UE&#39;s serving station identifier - list for all stations if not present |
| subscribe | [bool](#bool) |  | subscribe indicates whether to subscribe to events (e.g. ADD, UPDATE, and REMOVE) that occur after all uelinks have been streamed to the client |
| noReplay | [bool](#bool) |  | noReplay - do not replay the list of UELinks from before the request Used with subscribe, to only get new changes |
| noimsi | [bool](#bool) |  | noimsi - imsi is not needed in the response |





 


<a name="interface.c1.C1MessageType"></a>

### C1MessageType


| Name | Number | Description |
| ---- | ------ | ----------- |
| C1_MESSAGE_UNKNOWN | 0 |  |
| C1_GET_RNIBCELLS | 1 |  |
| C1_GET_RNIBUES | 2 |  |
| C1_GET_RNIBLINKS | 3 |  |
| C1_POST_HANDOVERS | 4 |  |
| C1_POST_RMMCONFIGURATION | 5 |  |



<a name="interface.c1.C1RNIBType"></a>

### C1RNIBType


| Name | Number | Description |
| ---- | ------ | ----------- |
| C1_RNIB_UNKNOWN | 0 |  |
| C1_RNIB_CELL | 1 |  |
| C1_RNIB_UE | 2 |  |
| C1_RNIB_LINK | 3 |  |



<a name="interface.c1.C1XICICPA"></a>

### C1XICICPA


| Name | Number | Description |
| ---- | ------ | ----------- |
| C1_XICIC_PA_DB_MINUS6 | 0 |  |
| C1_XICIC_PA_DB_MINUX4DOT77 | 1 |  |
| C1_XICIC_PA_DB_MINUS3 | 2 |  |
| C1_XICIC_PA_DB_MINUS1DOT77 | 3 |  |
| C1_XICIC_PA_DB_0 | 4 |  |
| C1_XICIC_PA_DB_1 | 5 |  |
| C1_XICIC_PA_DB_2 | 6 |  |
| C1_XICIC_PA_DB_3 | 7 |  |



<a name="interface.c1.StationPowerOffset"></a>

### StationPowerOffset
Enumerated Power offset - It is defined in E2 interface

| Name | Number | Description |
| ---- | ------ | ----------- |
| PA_DB_MINUS6 | 0 | TX power of the station - 6 dB |
| PA_DB_MINUX4DOT77 | 1 | TX power of the station - 4.77 dB |
| PA_DB_MINUS3 | 2 | TX power of the station - 3 dB |
| PA_DB_MINUS1DOT77 | 3 | TX power of the station - 1.77 dB |
| PA_DB_0 | 4 | TX power of the station - 0 dB |
| PA_DB_1 | 5 | TX power of the station &#43; 1 dB |
| PA_DB_2 | 6 | TX power of the station &#43; 2 dB |
| PA_DB_3 | 7 | TX power of the station &#43; 3 dB |


 

 


<a name="interface.c1.C1InterfaceService"></a>

### C1InterfaceService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListStations | [StationListRequest](#interface.c1.StationListRequest) | [StationInfo](#interface.c1.StationInfo) stream | ListStations returns a stream of base station records. |
| ListStationLinks | [StationLinkListRequest](#interface.c1.StationLinkListRequest) | [StationLinkInfo](#interface.c1.StationLinkInfo) stream | ListStationLinks returns a stream of links between neighboring base stations. |
| ListUELinks | [UELinkListRequest](#interface.c1.UELinkListRequest) | [UELinkInfo](#interface.c1.UELinkInfo) stream | ListUELinks returns a stream of UI and base station links; one-time or (later) continuous subscribe. |
| TriggerHandOver | [HandOverRequest](#interface.c1.HandOverRequest) | [HandOverResponse](#interface.c1.HandOverResponse) | TriggerHandOver returns a hand-over response indicating success or failure. HO app will ask to handover to multiple UEs |
| TriggerHandOverStream | [HandOverRequest](#interface.c1.HandOverRequest) stream | [HandOverResponse](#interface.c1.HandOverResponse) | TriggerHandOverStream is a version that stays open all the time and sends HO messages as they happen |
| SetRadioPower | [RadioPowerRequest](#interface.c1.RadioPowerRequest) | [RadioPowerResponse](#interface.c1.RadioPowerResponse) | SetRadioPower returns a response indicating success or failure. MLB app will ask to change transmission power to multiple stations |

 



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

