# Load Balancer Application

**onos-ric-mlb** is an application that adjusts the power level of cell towers, in inverse
proportion to the number of UEs attached. In this way, the application can influence UEs
to attach to different cell towers than their nearest one, and so can help balance the
load across cell towers.

It accesses onos-ric through its [C1](api/c1-interface.md) streaming gRPC interface.

It is packaged as a Cloud Native micro service, to be [deployed](deployment.md) on [Kubernetes].

[Kubernetes]: https://kubernetes.io/
[gRPC]: https://grpc.io/
