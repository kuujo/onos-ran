# onos-ran
RAN subsystem for ONOS (µONOS Architecture)

Presently, this is just a skeletal project established to provide support for the
Feb 2020 MWC demonstration.
However, we expected to start evolving quickly towards the long-term architecture.

## Interfaces
`onos-ran` presents a gRPC Northbound interface **c1**
and
it relies on a gRPC southbound interface **e2** which is provided by `ran-simulator`.

## Running
`onos-ran` can be run as:

* a standalone application `onos-ran --simulator <ip-address:port>`
* alongside `ran-simulator` and `sd-ran-gui` by using [docker-compose](../../ran-simulator/docs/docker-compose.md)
* in Kubernetes when deployed by a [Helm Chart](deployment.md)

# Applications
When running, `onos-ran` is accessed through its **c1** Interface. Future
applications include a Handover application and a MLB app. 

The [onos-cli](https://docs.onosproject.org/onos-cli/docs/setup/) application has
been extended to access some of the services provided by this interface

```bash
> go run github.com/onosproject/onos-cli/cmd/onos ran get --help
Get RAN resources

Usage:
  onos ran get [command]

Available Commands:
  stationlinks Get Station Links
  stations     Get Stations
  uelinks      Get UE Links

Flags:
  -h, --help   help for get

Global Flags:
      --no-tls                   if present, do not use TLS
      --service-address string   the onos-ran service address (default "onos-ran:5150")
      --tls-cert-path string     the path to the TLS certificate
      --tls-key-path string      the path to the TLS key

Use "onos ran get [command] --help" for more information about a command.
```
