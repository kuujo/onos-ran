# Near real time RAN Intelligent Controller (RIC)

**onos-ric** is the near real time RAN Intelligent Controller subsystem for ONOS (µONOS Architecture)
following the [O-RAN] Architecture.

It is built as a [Cloud Native] micro service to be deployed on [Kubernetes], and provides a C1 interface northbound and an
E2 interface southbound.

> It is a work in progress and follows developments in the O-RAN working groups.

## Interfaces
**onos-ric** presents a streaming [gRPC] Northbound interface [**c1**](api/c1-interface.md)
which is accessed by apps such as:

* onos-ran-ho - a [Handover](handover.md) application packaged as a Cloud Native micro service
* onos-ran-mlb - a [Load Balancing](loadbalancer.md) application packaged as a Cloud Native micro service
* onos-cli - the [command line](https://docs.onosproject.org/onos-cli/docs/cli/onos/) interface of µONOS
* onos-gui - the [Web](https://docs.onosproject.org/onos-gui/docs/ran-gui/) based interface of µONOS

On the southbound it relies on gRPC connections to:

* An interface [**e2ap**](api/e2ap.md) for [O-RAN] Access Protocol provided by `ran-simulator`.
* An interface [**e2sm**](api/e2sm.md) for [O-RAN] Service Model provided by `ran-simulator`.
* An interface [**e2**](api/e2-interface.md) for older X-RAN commands provided by `ran-simulator`.
* onos-topo - the µONOS [Topology Service](https://docs.onosproject.org/onos-topo/docs/api/device/)

**onos-ric** also gathers metrics using the [Prometheus] API, which can be used to track KPIs during operation.

## Running
`onos-ric` can be run only as a micro service in Kubernetes and is deployed by a [Helm] [Chart](deployment.md)

[O-RAN]: https://www.o-ran.org/
[Kubernetes]: https://kubernetes.io/
[Helm]: https://helm.sh/
[gRPC]: https://grpc.io/
[Cloud Native]: https://github.com/cncf/toc/blob/master/DEFINITION.md
[Prometheus]: https://prometheus.io/
