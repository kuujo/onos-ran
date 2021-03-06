# Deploying onos-ran

This guide deploys `onos-ran` through it's [Helm] chart assumes you have a
[Kubernetes] cluster running deployed in a namespace.

`onos-ran` Helm chart is based on Helm 3.0 version, with no need for the Tiller pod to be present. 
If you don't have a cluster running and want to try on your local machine please follow first 
the [Kubernetes] setup steps outlined in [deploy with Helm](https://docs.onosproject.org/developers/deploy_with_helm/).
The following steps assume you have the setup outlined in that page, including the `micro-onos` namespace configured. 

## Installing the Chart
To install the chart in the `micro-onos` namespace run from the root directory of
the `onos-helm-charts` repo the command:
```bash
helm install -n micro-onos onos-ran onos-ran
```
The output should be:
```bash
NAME: onos-ran
LAST DEPLOYED: Tue Feb  4 08:02:57 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

`helm install` assigns a unique name to the chart and displays all the k8s resources that were
created by it. To list the charts that are installed and view their statuses, run `helm ls`:

```bash
> helm ls
NAME         	NAMESPACE 	REVISION	UPDATED                                	STATUS  	CHART              	APP VERSION
onos-cli     	micro-onos	1       	2020-02-04 08:01:54.860813386 +0000 UTC	deployed	onos-cli-0.0.1     	1          
onos-ran     	micro-onos	1       	2020-02-04 08:02:17.663782372 +0000 UTC	deployed	onos-ran-0.0.1     	1          
ran-simulator	micro-onos	1       	2020-02-04 09:32:21.533299519 +0000 UTC	deployed	ran-simulator-0.0.1	1          
sd-ran-gui   	micro-onos	1       	2020-02-04 09:32:49.018099586 +0000 UTC	deployed	sd-ran-gui-0.0.1   	1  
```

> Here the service is shown running alongside `onos-cli`, `onos-ran` and the `sd-ran-gui`
> as these are usually deployed together to give a demo scenario. See the individual
> deployment instructions for these services.

To check whether the service is running correctly use `kubectl`
```bash
> kubectl get pods -n micro-onos
NAME                             READY   STATUS             RESTARTS   AGE
onos-cli-68bbf4f674-ssjt4        1/1     Running            0          18m
onos-ran-5fb8c6bdd7-xmcmq        1/1     Running            0          18m
ran-simulator-6f577597d8-5lcv8   1/1     Running            0          82s
sd-ran-gui-76ff54d85-fh72j       2/2     Running            0          54s
```

See Troubleshooting below if the `Status` is not `Running`

### Installing the chart in a different namespace.

Issue the `helm install` command substituting `micro-onos` with your namespace.
```bash
helm install -n <your_name_space> onos-ran onos-ran
```

### Troubleshoot
If your chart does not install or the pod is not running for some reason and/or you modified values Helm offers two flags to help you
debug your chart:  

* `--dry-run` check the chart without actually installing the pod. 
* `--debug` prints out more information about your chart

```bash
helm install -n micro-onos onos-ran --debug --dry-run onos-ran/
```

## Uninstalling the chart.

To remove the `onos-ran` pod issue
```bash
 helm delete -n micro-onos onos-ran
```

## Pod Information

To view the pods that are deployed, run `kubectl -n micro-onos get pods`.

[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[kind]: https://kind.sigs.k8s.io
