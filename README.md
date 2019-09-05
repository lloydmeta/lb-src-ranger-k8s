# lb-src-ranger

A k8s CRD that tends to the `loadBalancerSourceRanges` of your `LoadBalancer` services by sourcing the IPs from
a list of URLs that you specify, targeting your services via labels.

The URLs should return a list newline-separated CIDRs.

## Install

Install the CRD:

```shell
kubctl apply -f https://raw.githubusercontent.com/lloydmeta/lb-src-ranger-k8s/blob/lb-src-ranger.yaml
```

### See it in action

Optionally install the samples:

```shell
# Create a LoadBalancer service
kubctl apply -f https://raw.githubusercontent.com/lloydmeta/lb-src-ranger-k8s/blob/config/samples/dummy-service.yaml

# Create a LbSrcRanger
kubctl apply -f https://raw.githubusercontent.com/lloydmeta/lb-src-ranger-k8s/blob/config/samples/lbsrcranger_v1beta1_lbsrcranger.yaml
```

Get the service to see that its `loadBalancerSourceRanges` has been updated based on the URLs in the ranger.

## Dev Requirements

1. `kubectl` + `kubernetes`
2. `kustomize`