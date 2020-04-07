# Handover Application

**onos-ric-ho** is an application that calculates the need for the "handover" of a
UE to a different cell tower to the current serving tower.

It access onos-ric through its [C1](api/c1-interface.md) streaming [gRPC] interface.

It is packaged as a Cloud Native micro service, to be [deployed](deployment.md) on [Kubernetes].

[Kubernetes]: https://kubernetes.io/
[gRPC]: https://grpc.io/
