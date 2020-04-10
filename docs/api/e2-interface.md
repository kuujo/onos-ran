# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/sb/e2-interface.proto](#api/sb/e2-interface.proto)
    - [A1Param](#interface.e2.A1Param)
    - [A2Param](#interface.e2.A2Param)
    - [A3Param](#interface.e2.A3Param)
    - [A4Param](#interface.e2.A4Param)
    - [A5Param](#interface.e2.A5Param)
    - [A6Param](#interface.e2.A6Param)
    - [AddMeasId](#interface.e2.AddMeasId)
    - [BearerAdmissionRequest](#interface.e2.BearerAdmissionRequest)
    - [BearerAdmissionResponse](#interface.e2.BearerAdmissionResponse)
    - [BearerAdmissionStatus](#interface.e2.BearerAdmissionStatus)
    - [BearerReleaseInd](#interface.e2.BearerReleaseInd)
    - [CACap](#interface.e2.CACap)
    - [CandScell](#interface.e2.CandScell)
    - [CellConfigReport](#interface.e2.CellConfigReport)
    - [CellConfigRequest](#interface.e2.CellConfigRequest)
    - [DCCap](#interface.e2.DCCap)
    - [DelMeasId](#interface.e2.DelMeasId)
    - [ECGI](#interface.e2.ECGI)
    - [ERABParamsItem](#interface.e2.ERABParamsItem)
    - [ERABResponseItem](#interface.e2.ERABResponseItem)
    - [HOCause](#interface.e2.HOCause)
    - [HOComplete](#interface.e2.HOComplete)
    - [HOFailure](#interface.e2.HOFailure)
    - [HORequest](#interface.e2.HORequest)
    - [L2MeasConfig](#interface.e2.L2MeasConfig)
    - [L2ReportInterval](#interface.e2.L2ReportInterval)
    - [MeasCell](#interface.e2.MeasCell)
    - [MeasID](#interface.e2.MeasID)
    - [MeasIdAction](#interface.e2.MeasIdAction)
    - [MeasIdActionChoice](#interface.e2.MeasIdActionChoice)
    - [MeasObject](#interface.e2.MeasObject)
    - [PCIARFCN](#interface.e2.PCIARFCN)
    - [PDCPMeasReportPerUe](#interface.e2.PDCPMeasReportPerUe)
    - [PRBUsage](#interface.e2.PRBUsage)
    - [PerParam](#interface.e2.PerParam)
    - [PropScell](#interface.e2.PropScell)
    - [RRCMeasConfig](#interface.e2.RRCMeasConfig)
    - [RRMConfig](#interface.e2.RRMConfig)
    - [RRMConfigStatus](#interface.e2.RRMConfigStatus)
    - [RXSigReport](#interface.e2.RXSigReport)
    - [RadioMeasReportPerCell](#interface.e2.RadioMeasReportPerCell)
    - [RadioMeasReportPerUE](#interface.e2.RadioMeasReportPerUE)
    - [RadioRepPerServCell](#interface.e2.RadioRepPerServCell)
    - [ReportConfig](#interface.e2.ReportConfig)
    - [ReportParam](#interface.e2.ReportParam)
    - [ReportParamChoice](#interface.e2.ReportParamChoice)
    - [RxSigMeasReport](#interface.e2.RxSigMeasReport)
    - [ScellAdd](#interface.e2.ScellAdd)
    - [ScellAddStatus](#interface.e2.ScellAddStatus)
    - [ScellDelete](#interface.e2.ScellDelete)
    - [SchedMeasRepPerServCell](#interface.e2.SchedMeasRepPerServCell)
    - [SchedMeasReportPerCell](#interface.e2.SchedMeasReportPerCell)
    - [SchedMeasReportPerUE](#interface.e2.SchedMeasReportPerUE)
    - [SeNBAdd](#interface.e2.SeNBAdd)
    - [SeNBAddStatus](#interface.e2.SeNBAddStatus)
    - [SeNBDelete](#interface.e2.SeNBDelete)
    - [ThreasholdEUTRA](#interface.e2.ThreasholdEUTRA)
    - [ThresholdEUTRAChoice](#interface.e2.ThresholdEUTRAChoice)
    - [TrafficSplitConfig](#interface.e2.TrafficSplitConfig)
    - [TrafficSplitPercentage](#interface.e2.TrafficSplitPercentage)
    - [UEAMBR](#interface.e2.UEAMBR)
    - [UEAdmissionRequest](#interface.e2.UEAdmissionRequest)
    - [UEAdmissionResponse](#interface.e2.UEAdmissionResponse)
    - [UEAdmissionStatus](#interface.e2.UEAdmissionStatus)
    - [UECapabilityEnquiry](#interface.e2.UECapabilityEnquiry)
    - [UECapabilityInfo](#interface.e2.UECapabilityInfo)
    - [UEContextUpdate](#interface.e2.UEContextUpdate)
    - [UEReconfigInd](#interface.e2.UEReconfigInd)
    - [UEReleaseInd](#interface.e2.UEReleaseInd)
  
    - [AdmEstCause](#interface.e2.AdmEstCause)
    - [CACapClassDl](#interface.e2.CACapClassDl)
    - [CACapClassUl](#interface.e2.CACapClassUl)
    - [CADirection](#interface.e2.CADirection)
    - [DCCapDrbType](#interface.e2.DCCapDrbType)
    - [ERABDirection](#interface.e2.ERABDirection)
    - [ERABType](#interface.e2.ERABType)
    - [L2MeasReportIntervals](#interface.e2.L2MeasReportIntervals)
    - [MeasIdActionPR](#interface.e2.MeasIdActionPR)
    - [MessageType](#interface.e2.MessageType)
    - [PerParamReportIntervalMs](#interface.e2.PerParamReportIntervalMs)
    - [ReconfigCause](#interface.e2.ReconfigCause)
    - [ReleaseCause](#interface.e2.ReleaseCause)
    - [ReportParamPR](#interface.e2.ReportParamPR)
    - [ReportQuality](#interface.e2.ReportQuality)
    - [SuccessOrFailure](#interface.e2.SuccessOrFailure)
    - [ThresholdEUTRAPR](#interface.e2.ThresholdEUTRAPR)
    - [TimeToTrigger](#interface.e2.TimeToTrigger)
    - [TriggerQuantity](#interface.e2.TriggerQuantity)
    - [XICICPA](#interface.e2.XICICPA)
  
  
  

- [api/sb/e2ap.proto](#api/sb/e2ap.proto)
    - [E2apProtocolIE](#interface.e2ap.E2apProtocolIE)
    - [E2apProtocolIEsPair](#interface.e2ap.E2apProtocolIEsPair)
    - [ProtocolIEContainer](#interface.e2ap.ProtocolIEContainer)
    - [ProtocolIEField](#interface.e2ap.ProtocolIEField)
    - [ProtocolIESingleContainer](#interface.e2ap.ProtocolIESingleContainer)
  
    - [Criticality](#interface.e2ap.Criticality)
    - [Presence](#interface.e2ap.Presence)
    - [ProcedureCode](#interface.e2ap.ProcedureCode)
    - [ProtocolIEId](#interface.e2ap.ProtocolIEId)
    - [TriggeringMessage](#interface.e2ap.TriggeringMessage)
  
  
  

- [Scalar Value Types](#scalar-value-types)



<a name="api/sb/e2-interface.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/sb/e2-interface.proto
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


<a name="interface.e2.A1Param"></a>

### A1Param



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| a1Threshold | [ThreasholdEUTRA](#interface.e2.ThreasholdEUTRA) |  |  |






<a name="interface.e2.A2Param"></a>

### A2Param



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| a2Threshold | [ThreasholdEUTRA](#interface.e2.ThreasholdEUTRA) |  |  |






<a name="interface.e2.A3Param"></a>

### A3Param



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| a3Offset | [string](#string) |  |  |






<a name="interface.e2.A4Param"></a>

### A4Param



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| a4Threshold | [ThreasholdEUTRA](#interface.e2.ThreasholdEUTRA) |  |  |






<a name="interface.e2.A5Param"></a>

### A5Param



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| a5Threshold1 | [ThreasholdEUTRA](#interface.e2.ThreasholdEUTRA) |  |  |
| a5Threshold2 | [ThreasholdEUTRA](#interface.e2.ThreasholdEUTRA) |  |  |






<a name="interface.e2.A6Param"></a>

### A6Param



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| a6Offset | [string](#string) |  |  |






<a name="interface.e2.AddMeasId"></a>

### AddMeasId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| addMeasId | [string](#string) | repeated |  |






<a name="interface.e2.BearerAdmissionRequest"></a>

### BearerAdmissionRequest
BearerAdmissionRequest message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| ueAmbr | [UEAMBR](#interface.e2.UEAMBR) |  |  |
| numErabsList | [uint32](#uint32) |  |  |
| erabsParams | [ERABParamsItem](#interface.e2.ERABParamsItem) | repeated |  |






<a name="interface.e2.BearerAdmissionResponse"></a>

### BearerAdmissionResponse
BearerAdmissionResponse message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| numErabsList | [uint32](#uint32) |  |  |
| erabResponse | [ERABResponseItem](#interface.e2.ERABResponseItem) | repeated |  |






<a name="interface.e2.BearerAdmissionStatus"></a>

### BearerAdmissionStatus
BearerAdmissionStatus message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| numErabsList | [uint32](#uint32) |  |  |
| erabStatus | [ERABResponseItem](#interface.e2.ERABResponseItem) | repeated |  |






<a name="interface.e2.BearerReleaseInd"></a>

### BearerReleaseInd
BearerReleaseInd message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| numErabsList | [uint32](#uint32) |  |  |
| erabIds | [string](#string) | repeated |  |






<a name="interface.e2.CACap"></a>

### CACap



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| band | [string](#string) |  |  |
| caclassdl | [CACapClassDl](#interface.e2.CACapClassDl) |  |  |
| caclassul | [CACapClassUl](#interface.e2.CACapClassUl) |  |  |
| crossCarrierSched | [string](#string) |  |  |






<a name="interface.e2.CandScell"></a>

### CandScell



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pci | [uint32](#uint32) |  |  |
| earfcnDl | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |






<a name="interface.e2.CellConfigReport"></a>

### CellConfigReport
CellConfigReport message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| pci | [uint32](#uint32) |  |  |
| candScells | [CandScell](#interface.e2.CandScell) | repeated |  |
| earfcnDl | [string](#string) |  |  |
| earfcnUl | [string](#string) |  |  |
| rbsPerTtiDl | [uint32](#uint32) |  |  |
| rbsPerTtiUl | [uint32](#uint32) |  |  |
| numTxAntenna | [uint32](#uint32) |  |  |
| duplexMode | [string](#string) |  |  |
| maxNumConnectedUes | [uint32](#uint32) |  |  |
| maxNumConnectedBearers | [uint32](#uint32) |  |  |
| maxNumUesSchedPerTtiDl | [uint32](#uint32) |  |  |
| maxNumUesSchedPerTtiUl | [uint32](#uint32) |  |  |
| dlfsSchedEnable | [string](#string) |  |  |






<a name="interface.e2.CellConfigRequest"></a>

### CellConfigRequest
CellConfigRequest message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |






<a name="interface.e2.DCCap"></a>

### DCCap



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| drbTypeSplit | [DCCapDrbType](#interface.e2.DCCapDrbType) |  |  |






<a name="interface.e2.DelMeasId"></a>

### DelMeasId



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delMeasId | [string](#string) | repeated |  |






<a name="interface.e2.ECGI"></a>

### ECGI



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| plmnId | [string](#string) |  |  |
| ecid | [string](#string) |  |  |






<a name="interface.e2.ERABParamsItem"></a>

### ERABParamsItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| direction | [ERABDirection](#interface.e2.ERABDirection) |  |  |
| type | [ERABType](#interface.e2.ERABType) |  |  |
| qci | [string](#string) |  |  |
| arp | [uint32](#uint32) |  |  |
| gbrDl | [string](#string) |  |  |
| gbrUl | [string](#string) |  |  |
| mbrDl | [string](#string) |  |  |
| mbrUl | [string](#string) |  |  |






<a name="interface.e2.ERABResponseItem"></a>

### ERABResponseItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| decision | [SuccessOrFailure](#interface.e2.SuccessOrFailure) |  |  |






<a name="interface.e2.HOCause"></a>

### HOCause
HOCause message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgiS | [ECGI](#interface.e2.ECGI) |  |  |
| ecgiT | [ECGI](#interface.e2.ECGI) |  |  |
| hoCause | [string](#string) |  |  |
| hoTrigger | [RXSigReport](#interface.e2.RXSigReport) | repeated |  |






<a name="interface.e2.HOComplete"></a>

### HOComplete
HOComplete message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crntiNew | [string](#string) |  |  |
| ecgiS | [ECGI](#interface.e2.ECGI) |  |  |
| ecgiT | [ECGI](#interface.e2.ECGI) |  |  |






<a name="interface.e2.HOFailure"></a>

### HOFailure
HOFailure message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgiS | [ECGI](#interface.e2.ECGI) |  |  |
| hoFailureCause | [string](#string) |  | Not defined in standard yet |






<a name="interface.e2.HORequest"></a>

### HORequest
HORequest message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgiS | [ECGI](#interface.e2.ECGI) |  |  |
| ecgiT | [ECGI](#interface.e2.ECGI) |  |  |






<a name="interface.e2.L2MeasConfig"></a>

### L2MeasConfig
L2MeasConfig message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crnti | [string](#string) | repeated |  |
| reportIntervals | [L2ReportInterval](#interface.e2.L2ReportInterval) |  |  |






<a name="interface.e2.L2ReportInterval"></a>

### L2ReportInterval



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tRadioMeasReportPerUe | [L2MeasReportIntervals](#interface.e2.L2MeasReportIntervals) |  |  |
| tRadioMeasReportPerCell | [L2MeasReportIntervals](#interface.e2.L2MeasReportIntervals) |  |  |
| tSchedMeasReportPerUe | [L2MeasReportIntervals](#interface.e2.L2MeasReportIntervals) |  |  |
| tSchedMeasReportPerCell | [L2MeasReportIntervals](#interface.e2.L2MeasReportIntervals) |  |  |
| tPdcpMeasReportPerUe | [L2MeasReportIntervals](#interface.e2.L2MeasReportIntervals) |  |  |






<a name="interface.e2.MeasCell"></a>

### MeasCell



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pci | [uint32](#uint32) |  |  |
| cellIndividualOffset | [string](#string) |  |  |






<a name="interface.e2.MeasID"></a>

### MeasID



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| measObjectId | [string](#string) |  |  |
| reportConfigId | [string](#string) |  |  |
| action | [MeasIdAction](#interface.e2.MeasIdAction) |  |  |






<a name="interface.e2.MeasIdAction"></a>

### MeasIdAction



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| present | [MeasIdActionPR](#interface.e2.MeasIdActionPR) |  |  |
| choice | [MeasIdActionChoice](#interface.e2.MeasIdActionChoice) |  |  |






<a name="interface.e2.MeasIdActionChoice"></a>

### MeasIdActionChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| addMeasId | [AddMeasId](#interface.e2.AddMeasId) |  |  |
| delMeasId | [DelMeasId](#interface.e2.DelMeasId) |  |  |
| hototarget | [string](#string) |  |  |






<a name="interface.e2.MeasObject"></a>

### MeasObject



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dlFreq | [string](#string) |  |  |
| measCells | [MeasCell](#interface.e2.MeasCell) | repeated |  |






<a name="interface.e2.PCIARFCN"></a>

### PCIARFCN



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pci | [uint32](#uint32) |  |  |
| earfcnDl | [string](#string) |  |  |






<a name="interface.e2.PDCPMeasReportPerUe"></a>

### PDCPMeasReportPerUe
PdcpMeasReportPerUE message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crnti | [string](#string) |  |  |
| qciVals | [string](#string) | repeated |  |
| dataVolDl | [uint32](#uint32) | repeated |  |
| dataVolUl | [uint32](#uint32) | repeated |  |
| pktDelayDl | [uint32](#uint32) | repeated |  |
| pktDiscardRateDl | [uint32](#uint32) | repeated |  |
| pktLossRateDl | [uint32](#uint32) | repeated |  |
| pktLossRateUl | [uint32](#uint32) | repeated |  |
| throughputDl | [uint32](#uint32) | repeated |  |
| throughputUl | [uint32](#uint32) | repeated |  |






<a name="interface.e2.PRBUsage"></a>

### PRBUsage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| prbUsageDl | [string](#string) | repeated |  |
| prbUsageUl | [string](#string) | repeated |  |






<a name="interface.e2.PerParam"></a>

### PerParam



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reportIntervalMs | [PerParamReportIntervalMs](#interface.e2.PerParamReportIntervalMs) |  |  |






<a name="interface.e2.PropScell"></a>

### PropScell



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pciArfcn | [PCIARFCN](#interface.e2.PCIARFCN) |  |  |
| crossCarrierSchedEnable | [string](#string) |  |  |
| caDirection | [CADirection](#interface.e2.CADirection) |  |  |
| deactTimer | [string](#string) |  |  |






<a name="interface.e2.RRCMeasConfig"></a>

### RRCMeasConfig
RRCMeasConfig message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crnti | [string](#string) | repeated |  |
| measObjects | [MeasObject](#interface.e2.MeasObject) | repeated |  |
| reportConfigs | [ReportConfig](#interface.e2.ReportConfig) | repeated |  |
| measIds | [MeasID](#interface.e2.MeasID) | repeated |  |






<a name="interface.e2.RRMConfig"></a>

### RRMConfig
RRMConfig message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crnti | [string](#string) | repeated |  |
| pciArfcn | [PCIARFCN](#interface.e2.PCIARFCN) |  |  |
| pA | [XICICPA](#interface.e2.XICICPA) | repeated |  |
| startPrbDl | [uint32](#uint32) | repeated |  |
| endPrbDl | [uint32](#uint32) | repeated |  |
| subFrameBitmaskDl | [string](#string) | repeated |  |
| p0UePusch | [uint32](#uint32) | repeated |  |
| startPrbUl | [uint32](#uint32) | repeated |  |
| endPrbUl | [uint32](#uint32) | repeated |  |
| subFrameBitmaskUl | [string](#string) | repeated |  |






<a name="interface.e2.RRMConfigStatus"></a>

### RRMConfigStatus
RRMConfigStatus message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crnti | [string](#string) | repeated |  |
| status | [SuccessOrFailure](#interface.e2.SuccessOrFailure) | repeated |  |






<a name="interface.e2.RXSigReport"></a>

### RXSigReport



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pciArfcn | [PCIARFCN](#interface.e2.PCIARFCN) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| rsrp | [string](#string) |  |  |
| rsrq | [string](#string) |  | string measId = 5; not defined in our SBI |






<a name="interface.e2.RadioMeasReportPerCell"></a>

### RadioMeasReportPerCell
RadioMeasReportPerCell message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| puschIntfPwrHist | [uint32](#uint32) | repeated |  |
| pucchIntfPowerHist | [uint32](#uint32) | repeated |  |






<a name="interface.e2.RadioMeasReportPerUE"></a>

### RadioMeasReportPerUE
RadioMeasReportPerUE message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crnti | [string](#string) |  |  |
| radioReportServCells | [RadioRepPerServCell](#interface.e2.RadioRepPerServCell) | repeated |  |






<a name="interface.e2.RadioRepPerServCell"></a>

### RadioRepPerServCell
L2MeasureReports (periodic) messages


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| cqiHist | [uint32](#uint32) | repeated |  |
| riHist | [string](#string) | repeated |  |
| puschSinrHist | [string](#string) | repeated |  |
| pucchSinrHist | [string](#string) | repeated |  |






<a name="interface.e2.ReportConfig"></a>

### ReportConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reportParams | [ReportParam](#interface.e2.ReportParam) |  |  |
| triggerQuantity | [TriggerQuantity](#interface.e2.TriggerQuantity) |  |  |
| reportQuality | [ReportQuality](#interface.e2.ReportQuality) |  |  |






<a name="interface.e2.ReportParam"></a>

### ReportParam



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| present | [ReportParamPR](#interface.e2.ReportParamPR) |  |  |
| choice | [ReportParamChoice](#interface.e2.ReportParamChoice) |  |  |
| hysteresis | [string](#string) |  |  |
| timetotrigger | [TimeToTrigger](#interface.e2.TimeToTrigger) |  |  |






<a name="interface.e2.ReportParamChoice"></a>

### ReportParamChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| perParam | [PerParam](#interface.e2.PerParam) |  |  |
| a1Param | [A1Param](#interface.e2.A1Param) |  |  |
| a2Param | [A2Param](#interface.e2.A2Param) |  |  |
| a3Param | [A3Param](#interface.e2.A3Param) |  |  |
| a4Param | [A4Param](#interface.e2.A4Param) |  |  |
| a5Param | [A5Param](#interface.e2.A5Param) |  |  |
| a6Param | [A6Param](#interface.e2.A6Param) |  |  |






<a name="interface.e2.RxSigMeasReport"></a>

### RxSigMeasReport
RxSigMeasReport message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) | repeated |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| cellMeasReport | [RXSigReport](#interface.e2.RXSigReport) | repeated |  |






<a name="interface.e2.ScellAdd"></a>

### ScellAdd
ScellAdd message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| scellsProp | [PropScell](#interface.e2.PropScell) | repeated |  |






<a name="interface.e2.ScellAddStatus"></a>

### ScellAddStatus
ScellAddStatus message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| scellsInd | [PCIARFCN](#interface.e2.PCIARFCN) | repeated |  |
| status | [SuccessOrFailure](#interface.e2.SuccessOrFailure) | repeated |  |






<a name="interface.e2.ScellDelete"></a>

### ScellDelete
ScellDelete message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| scellsInd | [PCIARFCN](#interface.e2.PCIARFCN) | repeated |  |






<a name="interface.e2.SchedMeasRepPerServCell"></a>

### SchedMeasRepPerServCell



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| qciVals | [string](#string) | repeated |  |
| prbUsage | [PRBUsage](#interface.e2.PRBUsage) |  |  |
| mcsDl | [string](#string) | repeated |  |
| numSchedTtisDl | [string](#string) | repeated |  |
| mcsUl | [string](#string) | repeated |  |
| numSchedTtisUl | [string](#string) | repeated |  |
| rankDl1 | [string](#string) | repeated |  |
| rankDl2 | [string](#string) | repeated |  |






<a name="interface.e2.SchedMeasReportPerCell"></a>

### SchedMeasReportPerCell
SchedMeasReportPerCell message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| qciVals | [string](#string) | repeated |  |
| prbUsagePcell | [PRBUsage](#interface.e2.PRBUsage) |  |  |
| prbUsageScell | [PRBUsage](#interface.e2.PRBUsage) |  |  |






<a name="interface.e2.SchedMeasReportPerUE"></a>

### SchedMeasReportPerUE
SchedMeasReportPerUE message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crnti | [string](#string) |  |  |
| schedReportServCells | [SchedMeasRepPerServCell](#interface.e2.SchedMeasRepPerServCell) | repeated |  |






<a name="interface.e2.SeNBAdd"></a>

### SeNBAdd
SeNBAdd message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| mEcgi | [ECGI](#interface.e2.ECGI) |  |  |
| sEcgi | [ECGI](#interface.e2.ECGI) |  |  |






<a name="interface.e2.SeNBAddStatus"></a>

### SeNBAddStatus
SeNBAddStatus message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| status | [SuccessOrFailure](#interface.e2.SuccessOrFailure) |  |  |






<a name="interface.e2.SeNBDelete"></a>

### SeNBDelete
SeNBDelete message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| mEcgi | [ECGI](#interface.e2.ECGI) |  |  |
| sEcgi | [ECGI](#interface.e2.ECGI) |  |  |






<a name="interface.e2.ThreasholdEUTRA"></a>

### ThreasholdEUTRA



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| present | [ThresholdEUTRAPR](#interface.e2.ThresholdEUTRAPR) |  |  |
| choice | [ThresholdEUTRAChoice](#interface.e2.ThresholdEUTRAChoice) |  |  |






<a name="interface.e2.ThresholdEUTRAChoice"></a>

### ThresholdEUTRAChoice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thresholdRSRP | [string](#string) |  |  |
| thresholdRSRQ | [string](#string) |  |  |






<a name="interface.e2.TrafficSplitConfig"></a>

### TrafficSplitConfig
TrafficSplitConfig message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| trafficSplitPercentage | [TrafficSplitPercentage](#interface.e2.TrafficSplitPercentage) | repeated |  |






<a name="interface.e2.TrafficSplitPercentage"></a>

### TrafficSplitPercentage



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| trafficPercentageDl | [string](#string) |  |  |
| trafficPercentageUl | [string](#string) |  |  |






<a name="interface.e2.UEAMBR"></a>

### UEAMBR



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ambrDl | [string](#string) |  |  |
| ambrUl | [string](#string) |  |  |






<a name="interface.e2.UEAdmissionRequest"></a>

### UEAdmissionRequest
UEAdmissionRequest message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| admissionEstCause | [AdmEstCause](#interface.e2.AdmEstCause) |  |  |
| imsi | [uint64](#uint64) |  |  |






<a name="interface.e2.UEAdmissionResponse"></a>

### UEAdmissionResponse
UEAdmissionResponse message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| admissionEstResponse | [SuccessOrFailure](#interface.e2.SuccessOrFailure) |  |  |






<a name="interface.e2.UEAdmissionStatus"></a>

### UEAdmissionStatus
UEAdmissionStatus message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| admissionEstStatus | [SuccessOrFailure](#interface.e2.SuccessOrFailure) |  |  |






<a name="interface.e2.UECapabilityEnquiry"></a>

### UECapabilityEnquiry
UECapabilityEnqyuiry message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |






<a name="interface.e2.UECapabilityInfo"></a>

### UECapabilityInfo
UECapabilityInfo message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| caCap | [CACap](#interface.e2.CACap) |  |  |
| dcCap | [DCCap](#interface.e2.DCCap) |  |  |






<a name="interface.e2.UEContextUpdate"></a>

### UEContextUpdate
UEContextUpdate message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| mmeUeS1apId | [string](#string) |  |  |
| enbUeS1apId | [string](#string) |  |  |
| imsi | [string](#string) |  |  |






<a name="interface.e2.UEReconfigInd"></a>

### UEReconfigInd
UEReconfigInd message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crntiOld | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| crntiNew | [string](#string) |  |  |
| reconfigCause | [ReconfigCause](#interface.e2.ReconfigCause) |  |  |






<a name="interface.e2.UEReleaseInd"></a>

### UEReleaseInd
UEReleaseInd message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| crnti | [string](#string) |  |  |
| ecgi | [ECGI](#interface.e2.ECGI) |  |  |
| releaseCause | [ReleaseCause](#interface.e2.ReleaseCause) |  |  |





 


<a name="interface.e2.AdmEstCause"></a>

### AdmEstCause


| Name | Number | Description |
| ---- | ------ | ----------- |
| EMERGENCY | 0 |  |
| HIGHHP_ACCESS | 1 |  |
| MT_ACCESS | 2 |  |
| MO_SIGNALLING | 3 |  |
| MO_DATA | 4 |  |



<a name="interface.e2.CACapClassDl"></a>

### CACapClassDl


| Name | Number | Description |
| ---- | ------ | ----------- |
| CACAP_CLASSDL_A | 0 |  |
| CACAP_CLASSDL_B | 1 |  |
| CACAP_CLASSDL_C | 2 |  |
| CACAP_CLASSDL_D | 3 |  |
| CACAP_CLASSDL_E | 4 |  |
| CACAP_CLASSDL_F | 5 |  |



<a name="interface.e2.CACapClassUl"></a>

### CACapClassUl


| Name | Number | Description |
| ---- | ------ | ----------- |
| CACAP_CLASSUL_A | 0 |  |
| CACAP_CLASSUL_B | 1 |  |
| CACAP_CLASSUL_C | 2 |  |
| CACAP_CLASSUL_D | 3 |  |
| CACAP_CLASSUL_E | 4 |  |
| CACAP_CLASSUL_F | 5 |  |



<a name="interface.e2.CADirection"></a>

### CADirection


| Name | Number | Description |
| ---- | ------ | ----------- |
| CADIRECTION_DL | 0 |  |
| CADIRECTION_UL | 1 |  |
| CADIRECTION_BOTH | 2 |  |



<a name="interface.e2.DCCapDrbType"></a>

### DCCapDrbType


| Name | Number | Description |
| ---- | ------ | ----------- |
| DCCAP_DRBTYPE_SUPPORTED | 0 |  |



<a name="interface.e2.ERABDirection"></a>

### ERABDirection


| Name | Number | Description |
| ---- | ------ | ----------- |
| DL | 0 |  |
| UL | 1 |  |
| BOTH | 2 |  |



<a name="interface.e2.ERABType"></a>

### ERABType


| Name | Number | Description |
| ---- | ------ | ----------- |
| ERAB_DEFAULT | 0 |  |
| ERAB_DEDICATED | 1 |  |



<a name="interface.e2.L2MeasReportIntervals"></a>

### L2MeasReportIntervals


| Name | Number | Description |
| ---- | ------ | ----------- |
| L2_MEAS_REPORT_INTERVAL_NO_REPORT | 0 |  |
| L2_MEAS_REPORT_INTERVAL_MS_10 | 1 |  |
| L2_MEAS_REPORT_INTERVAL_MS_50 | 2 |  |
| L2_MEAS_REPORT_INTERVAL_MS_100 | 3 |  |
| L2_MEAS_REPORT_INTERVAL_MS_200 | 4 |  |
| L2_MEAS_REPORT_INTERVAL_MS_500 | 5 |  |
| L2_MEAS_REPORT_INTERVAL_MS_1024 | 6 |  |
| L2_MEAS_REPORT_INTERVAL_S_10 | 7 |  |
| L2_MEAS_REPORT_INTERVAL_MIN_1 | 8 |  |
| L2_MEAS_REPORT_INTERVAL_MIN_2 | 9 |  |
| L2_MEAS_REPORT_INTERVAL_MIN_5 | 10 |  |



<a name="interface.e2.MeasIdActionPR"></a>

### MeasIdActionPR


| Name | Number | Description |
| ---- | ------ | ----------- |
| MEAS_ID_ACTION_PR_NOTHING | 0 |  |
| MEAS_ID_ACTION_PR_ADDMEASID | 1 |  |
| MEAS_ID_ACTION_PR_DELMEASID | 2 |  |
| MEAS_ID_ACTION_PR_HOTOTARGET | 3 |  |



<a name="interface.e2.MessageType"></a>

### MessageType
Definition of all message types in E2 interface

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN_MESSAGE | 0 |  |
| CELL_CONFIG_REQUEST | 1 | Defined, but not used yet |
| CELL_CONFIG_REPORT | 2 | Defined and used at xranc and ric-api-gw |
| UE_ADMISSION_REQUEST | 3 | Defined, but not used yet |
| UE_ADMISSION_RESPONSE | 4 | Defined, but not used yet |
| UE_ADMISSION_STATUS | 5 | Defined and used at xranc and ric-api-gw |
| UE_CONTEXT_UPDATE | 6 | Defined and used at xranc and ric-api-gw |
| UE_RECONFIG_IND | 7 | Defined, but not used yet |
| UE_RELEASE_IND | 8 | Defined, but not used yet |
| BEARER_ADMISSION_REQUEST | 9 | Defined, but not used yet |
| BEARER_ADMISSION_RESPONSE | 10 | Defined, but not used yet |
| BEARER_ADMISSION_STATUS | 11 | Defined, but not used yet |
| BEARER_RELEASE_IND | 12 | Defined, but not used yet |
| HO_REQUEST | 13 | Defined, but not used yet |
| HO_FAILURE | 14 | Defined, but not used yet |
| HO_COMPLETE | 15 | Defined, but not used yet |
| RXSIG_MEAS_REPORT | 16 | Defined, but not used yet |
| L2_MEAS_CONFIG | 17 | Defined, but not used yet |
| RADIO_MEAS_REPORT_PER_UE | 18 | Defined, but not used yet |
| RADIO_MEAS_REPORT_PER_CELL | 19 | Defined, but not used yet |
| SCHED_MEAS_REPORT_PER_UE | 20 | Defined, but not used yet |
| SCHED_MEAS_REPORT_PER_CELL | 21 | Defined, but not used yet |
| PDCP_MEAS_REPORT_PER_UE | 22 | Defined, but not used yet |
| UE_CAPABILITY_INFO | 23 | Defined, but not used yet |
| UE_CAPABILITY_ENQUIRY | 24 | Defined, but not used yet |
| SCELL_ADD | 25 | Defined, but not used yet |
| SCELL_ADD_STATUS | 26 | Defined, but not used yet |
| SCELL_DELETE | 27 | Defined, but not used yet |
| RRM_CONFIG | 28 | Defined, but not used yet |
| RRM_CONFIG_STATUS | 29 | Defined, but not used yet |
| SENB_ADD | 30 | Defined, but not used yet |
| SENB_ADD_STATUS | 31 | Defined, but not used yet |
| SENB_DELETE | 32 | Defined, but not used yet |
| TRAFFIC_SPLIT_CONFIG | 33 | Defined, but not used yet |
| HO_CAUSE | 34 | Defined, but not used yet |
| RRC_MEAS_CONFIG | 35 | Defined, but not used yet |



<a name="interface.e2.PerParamReportIntervalMs"></a>

### PerParamReportIntervalMs


| Name | Number | Description |
| ---- | ------ | ----------- |
| PER_PARAM_MS_120 | 0 |  |
| PER_PARAM_MS_240 | 1 |  |
| PER_PARAM_MS_480 | 2 |  |
| PER_PARAM_MS_640 | 3 |  |
| PER_PARAM_MS_1024 | 4 |  |
| PER_PARAM_MS_2048 | 5 |  |
| PER_PARAM_MS_5120 | 6 |  |
| PER_PARAM_MS_10240 | 7 |  |
| PER_PARAM_MIN_1 | 8 |  |
| PER_PARAM_MIN_6 | 9 |  |
| PER_PARAM_MIN_12 | 10 |  |
| PER_PARAM_MIN_30 | 11 |  |
| PER_PARAM_MIN_60 | 12 |  |



<a name="interface.e2.ReconfigCause"></a>

### ReconfigCause


| Name | Number | Description |
| ---- | ------ | ----------- |
| RECONFIG_RLF | 0 |  |
| RECONFIG_HO_FAIL | 1 |  |
| RECONFIG_OTHERS | 2 |  |



<a name="interface.e2.ReleaseCause"></a>

### ReleaseCause


| Name | Number | Description |
| ---- | ------ | ----------- |
| RELEASE_INACTIVITY | 0 |  |
| RELEASE_RLF | 1 |  |
| RELEASE_OTHERS | 2 |  |



<a name="interface.e2.ReportParamPR"></a>

### ReportParamPR


| Name | Number | Description |
| ---- | ------ | ----------- |
| REPORT_PARAM_PR_NOTHING | 0 |  |
| REPORT_PARAM_PR_PER_PARAM | 1 |  |
| REPORT_PARAM_PR_A1PARAM | 2 |  |
| REPORT_PARAM_PR_A2PARAM | 3 |  |
| REPORT_PARAM_PR_A3PARAM | 4 |  |
| REPORT_PARAM_PR_A4PARAM | 5 |  |
| REPORT_PARAM_PR_A5PARAM | 6 |  |
| REPORT_PARAM_PR_A6PARAM | 7 |  |



<a name="interface.e2.ReportQuality"></a>

### ReportQuality


| Name | Number | Description |
| ---- | ------ | ----------- |
| RQ_SAME | 0 |  |
| RQ_BOTH | 1 |  |



<a name="interface.e2.SuccessOrFailure"></a>

### SuccessOrFailure


| Name | Number | Description |
| ---- | ------ | ----------- |
| SUCCESS | 0 |  |
| FAILURE | 1 |  |



<a name="interface.e2.ThresholdEUTRAPR"></a>

### ThresholdEUTRAPR


| Name | Number | Description |
| ---- | ------ | ----------- |
| THRESHOLD_EUTRA_PR_NOTHING | 0 |  |
| THRESHOLD_EUTRA_PR_RSRP | 1 |  |
| THRESHOLD_EUTRA_PR_RSRQ | 2 |  |



<a name="interface.e2.TimeToTrigger"></a>

### TimeToTrigger


| Name | Number | Description |
| ---- | ------ | ----------- |
| TTT_MS0 | 0 |  |
| TTT_MS40 | 1 |  |
| TTT_MS64 | 2 |  |
| TTT_MS80 | 3 |  |
| TTT_MS100 | 4 |  |
| TTT_MS128 | 5 |  |
| TTT_MS160 | 6 |  |
| TTT_MS256 | 7 |  |
| TTT_MS320 | 8 |  |
| TTT_MS480 | 9 |  |
| TTT_MS512 | 10 |  |
| TTT_MS640 | 11 |  |
| TTT_MS1024 | 12 |  |
| TTT_MS1280 | 13 |  |
| TTT_MS2560 | 14 |  |
| TTT_MS5120 | 15 |  |



<a name="interface.e2.TriggerQuantity"></a>

### TriggerQuantity


| Name | Number | Description |
| ---- | ------ | ----------- |
| TQ_RSRP | 0 |  |
| TQ_RSRQ | 1 |  |



<a name="interface.e2.XICICPA"></a>

### XICICPA


| Name | Number | Description |
| ---- | ------ | ----------- |
| XICIC_PA_DB_MINUS6 | 0 |  |
| XICIC_PA_DB_MINUX4DOT77 | 1 |  |
| XICIC_PA_DB_MINUS3 | 2 |  |
| XICIC_PA_DB_MINUS1DOT77 | 3 |  |
| XICIC_PA_DB_0 | 4 |  |
| XICIC_PA_DB_1 | 5 |  |
| XICIC_PA_DB_2 | 6 |  |
| XICIC_PA_DB_3 | 7 |  |


 

 

 



<a name="api/sb/e2ap.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/sb/e2ap.proto
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


<a name="interface.e2ap.E2apProtocolIE"></a>

### E2apProtocolIE



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [ProtocolIEId](#interface.e2ap.ProtocolIEId) |  |  |
| criticality | [Criticality](#interface.e2ap.Criticality) |  |  |
| presence | [Presence](#interface.e2ap.Presence) |  | value |






<a name="interface.e2ap.E2apProtocolIEsPair"></a>

### E2apProtocolIEsPair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [ProtocolIEId](#interface.e2ap.ProtocolIEId) |  |  |
| firstCriticality | [Criticality](#interface.e2ap.Criticality) |  |  |
| secondCriticality | [Criticality](#interface.e2ap.Criticality) |  | firstValue = 3; |
| presence | [Presence](#interface.e2ap.Presence) |  | secondValue = 5; |






<a name="interface.e2ap.ProtocolIEContainer"></a>

### ProtocolIEContainer







<a name="interface.e2ap.ProtocolIEField"></a>

### ProtocolIEField







<a name="interface.e2ap.ProtocolIESingleContainer"></a>

### ProtocolIESingleContainer






 


<a name="interface.e2ap.Criticality"></a>

### Criticality


| Name | Number | Description |
| ---- | ------ | ----------- |
| REJECT | 0 |  |
| IGNORE | 1 |  |
| NOTIFY | 2 |  |



<a name="interface.e2ap.Presence"></a>

### Presence


| Name | Number | Description |
| ---- | ------ | ----------- |
| OPTIONAL | 0 |  |
| CONDITIONAL | 1 |  |
| MANDATORY | 2 |  |



<a name="interface.e2ap.ProcedureCode"></a>

### ProcedureCode


| Name | Number | Description |
| ---- | ------ | ----------- |
| PC_INVALID | 0 |  |
| E2_SETUP | 1 |  |
| ERROR_INDICATION | 2 |  |
| RESET | 3 |  |
| RIC_CONTROL | 4 |  |
| RIC_INDICATION | 5 |  |
| RIC_SERVICE_QUERY | 6 |  |
| RIC_SERVICE_UPDATE | 7 |  |
| RIC_SUBSCRIPTION | 8 |  |
| RIC_SUBSCRIPTION_DELETE | 9 |  |



<a name="interface.e2ap.ProtocolIEId"></a>

### ProtocolIEId


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNDEFINED | 0 |  |
| CAUSE | 1 |  |
| CRITICALITY_DIAGNOSTICS | 2 |  |
| GLOBAL_E2_NODE_ID | 3 |  |
| GLOBAL_RIC_ID | 4 |  |
| RAN_FUNCTION_ID | 5 |  |
| RAN_FUNCTION_ID_ITEM | 6 |  |
| RAN_FUNCTION_IE_CAUSE_ITEM | 7 |  |
| RAN_FUNCTION_ITEM | 8 |  |
| RAN_FUNCTIONS_ACCEPTED | 9 |  |
| RAN_FUNCTIONS_ADDED | 10 |  |
| RAN_FUNCTIONS_DELETED | 11 |  |
| RAN_FUNCTIONS_MODIFIED | 12 |  |
| RAN_FUNCTIONS_REJECTED | 13 |  |
| RIC_ACTION_ADMITTED_ITEM | 14 |  |
| RIC_ACTION_ID | 15 |  |
| RIC_ACTION_NOT_ADMITTED_ITEM | 16 |  |
| RIC_ACTIONS_ADMITTED | 17 |  |
| RIC_ACTIONS_NOT_ADMITTED | 18 |  |
| RIC_ACTION_TO_BE_SETUP_ITEM | 19 |  |
| RIC_CALL_PROCESS_ID | 20 |  |
| RIC_CONTROL_ACK_REQUEST | 21 |  |
| RIC_CONTROL_HEADER | 22 |  |
| RIC_CONTROL_MESSAGE | 23 |  |
| RIC_CONTROL_STATUS | 24 |  |
| RIC_INDICATION_HEADER | 25 |  |
| RIC_INDICATION_MESSAGE | 26 |  |
| RIC_INDICATION_SN | 27 |  |
| RIC_INDICATION_TYPE | 28 |  |
| RIC_REQUEST_ID | 29 |  |
| RIC_SUBSCRIPTION_DETAILS | 30 |  |
| TIME_TO_WAIT | 31 |  |
| RIC_CONTROL_OUTCOME | 32 |  |



<a name="interface.e2ap.TriggeringMessage"></a>

### TriggeringMessage


| Name | Number | Description |
| ---- | ------ | ----------- |
| INITIATING_MESSAGE | 0 |  |
| SUCCESSFUL_OUTCOME | 1 |  |
| UNSUCCESSFULL_OUTCOME | 2 |  |


 

 

 



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

