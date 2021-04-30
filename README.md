## KPNG - IPVS Backend

This is the PoC of an IPVS Backend, using cmd.Run instead of netlink

This is an early WIP and should not be used in anyway as production

The idea here is to deep dive into [kpng](https://github.com/kubernetes-sigs/kpng) 
work, and how this can be used with multiple backends

To run this:

* Start a kpng server with kubernetes access somewhere else
```
kpng kube to-api  --kubeconfig mykubeconfig.conf
```

* On other server (or the same), run this backend:
```
ipvs --nodeport-address 192.168.0.150
```

You need to have ipvsadm installed in the machine that will run kpng.
SCTP support only works on ipvsadm v1.30 or later

## Building the binary

Just use `make build`

## Building the container image

Just use `make image`

