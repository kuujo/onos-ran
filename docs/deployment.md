# Deploying onos-ric, onos-ric-ho and onos-ric-mlb individually

This guide deploys `onos-ric` through it's [Helm] chart assumes you have a
[Kubernetes] cluster running deployed in a namespace.

These Helm charts are based on Helm 3.0 version, with no need for the Tiller pod to be present.

> If you don't have a cluster running and want to try on your local machine
> please follow first the [Kubernetes] setup steps outlined in
> [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).

## SD-RAN Umbrella chart
For deployment of all of the µONOS services together in one step use the
[SD-RAN Umbrella chart](https://docs.onosproject.org/developers/deploy_with_helm/)

> The following steps assume you have the prerequisites outlined in that page,
> including the `micro-onos` namespace configured, but not the SD-RAN umbrella chart,
> and that you want to install services individually.

## Installing the Chart
To install the charts individually in the `micro-onos` namespace run from the root directory of
the `onos-helm-charts` repo the command:
```bash
helm install -n micro-onos onos-ric onos-ric
```
The output should be:
```bash
NAME: onos-ric
LAST DEPLOYED: Tue Feb  4 08:02:57 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

### Handover and Loadbalancer
Similarly the [Handover](handover.md) and [Load Balancer](loadbalancer.md)
applications can be deployed through Helm:
```bash
helm install -n micro-onos onos-ric-ho onos-ric-ho
helm install -n micro-onos onos-ric-mlb onos-ric-mlb
```

`helm install` assigns a unique name to the chart and displays all the k8s resources that were
created by it.

To list the charts that are installed and view their statuses, run `helm ls`:

```bash
> helm ls
NAME         	NAMESPACE 	REVISION	UPDATED                                	STATUS  	CHART              	APP VERSION
onos-cli     	micro-onos	1       	2020-02-04 08:01:54.860813386 +0000 UTC	deployed	onos-cli-0.0.1     	1
onos-ric     	micro-onos	1       	2020-02-04 08:02:17.663782372 +0000 UTC	deployed	onos-ric-0.0.1     	1
onos-ric-ho  	micro-onos	1       	2020-02-04 08:02:27.076256472 +0000 UTC	deployed	onos-ric-ho-0.0.1	1
onos-ric-mlb 	micro-onos	1       	2020-02-04 08:02:34.943652464 +0000 UTC	deployed	onos-ric-mlb-0.0.1	1
ran-simulator	micro-onos	1       	2020-02-04 09:32:21.533299519 +0000 UTC	deployed	ran-simulator-0.0.1	1
sd-ran-gui   	micro-onos	1       	2020-02-04 09:32:49.018099586 +0000 UTC	deployed	sd-ran-gui-0.0.1   	1
```

> Here the services are shown running alongside `onos-cli`, `ran-simulator` and the `onos-gui`
> as these are usually deployed together to give a demo scenario. See the individual
> deployment instructions for these services.

To check whether the service is running correctly use `kubectl`
```bash
> kubectl get pods -n micro-onos
NAME                             READY   STATUS             RESTARTS   AGE
onos-cli-68bbf4f674-ssjt4        1/1     Running            0          18m
onos-ric-5fb8c6bdd7-xmcmq        1/1     Running            0          18m
onos-ric-ho-9cb54ad7c2-defgh     1/1     Running            0          18m
onos-ric-mlb-46dcfbdc21-pqrst    1/1     Running            0          17m
ran-simulator-6f577597d8-5lcv8   1/1     Running            0          82s
onos-gui-76ff54d85-fh72j         2/2     Running            0          54s
```

See Troubleshooting below if the `Status` is not `Running`

### Installing the chart in a different namespace.

Issue the `helm install` command substituting `micro-onos` with your namespace.
```bash
helm install -n <your_name_space> onos-ric onos-ric
```

### Troubleshoot
If your chart does not install or the pod is not running for some reason and/or you modified values Helm offers two flags to help you
debug your chart:  

* `--dry-run` check the chart without actually installing the pod. 
* `--debug` prints out more information about your chart

```bash
helm install -n micro-onos onos-ric --debug --dry-run onos-ric/
```

## Uninstalling the chart.

To remove the `onos-ric` pod issue
```bash
 helm delete -n micro-onos onos-ric
```

## Pod Information

To view the pods that are deployed, run `kubectl -n micro-onos get pods`.

[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io
