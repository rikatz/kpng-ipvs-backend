package ipvs

import (
	"fmt"
	"os/exec"
	"strings"

	"k8s.io/klog"

	"sigs.k8s.io/kpng/pkg/api/localnetv1"
	"sigs.k8s.io/kpng/pkg/client"
)

const (
	// TODO: This is configurable :)
	DefaultAlgo   = "rr"
	DefaultWeight = 1
)

func Callback(ch <-chan *client.ServiceEndpoints) {

	var ipvsCfg strings.Builder
	var err error

	for serviceEndpoints := range ch {

		svc := serviceEndpoints.Service
		endpoints := serviceEndpoints.Endpoints

		//var svcRule strings.Builder
		// If this is a ClusterIP non headless
		if (svc.Type == "ClusterIP" && svc.GetIPs().ClusterIP != "None") || svc.Type == "NodePort" || svc.Type == "LoadBalancer" {
			// TODO: Verify if the IP Address exists on the dummy interface, otherwise it need to be
			// added
			// Future art: build some tree as Mikael did in nft, and verify if and where are the differences
			// to remove unused ClusterIP Addresses
			var nodePort bool
			if svc.Type == "NodePort" || svc.Type == "LoadBalancer" {
				nodePort = true
			}
			cip, err := buildClusterIP(svc, endpoints, nodePort)
			if err != nil {
				klog.Warningf("problem creating the service: %s", err)
			}
			ipvsCfg.WriteString(cip)
		}
	}
	fmt.Printf("%s", ipvsCfg.String())

	if OnlyOutput != nil && !*OnlyOutput {
		fmt.Println("Running clear")
		ipvsClear := exec.Command(*IPVSAdmPath, "--clear")

		err = ipvsClear.Run()
		if err != nil {
			fmt.Println("Error")
			klog.Errorf("failed to clear ipvs table: %s", err)
		}

		fmt.Println("Running restore")
		ipvsRestore := exec.Command(*IPVSAdmPath, "--restore")
		ipvsRestore.Stdin = strings.NewReader(ipvsCfg.String())

		err = ipvsRestore.Run()
		if err != nil {
			fmt.Println("Error")
			klog.Errorf("failed to execute ipvsadm restore: %s", err)
		}
	}
	ipvsCfg.Reset()
}

func buildClusterIP(svc *localnetv1.Service, eps []*localnetv1.Endpoint, nodePort bool) (string, error) {
	var svcString strings.Builder
	for _, port := range svc.Ports {
		var proto string
		switch port.Protocol {
		case localnetv1.Protocol_TCP:
			proto = "-t"
		case localnetv1.Protocol_SCTP:
			proto = "--sctp-service "
		case localnetv1.Protocol_UDP:
			proto = "-u"
		default:
			return "", fmt.Errorf("service %s/%s uses an unknown protocol", svc.Namespace, svc.Name)
		}
		ipPortAlgo := fmt.Sprintf("-A %s %s:%d -s %s\n", proto, svc.GetIPs().ClusterIP, port.GetTargetPort(), DefaultAlgo)
		svcString.WriteString(ipPortAlgo)

		endpoints := buildEndponts(svc.GetIPs().ClusterIP, proto, port.GetTargetPort(), port.GetPort(), eps)
		svcString.WriteString(endpoints)

		if nodePort {
			for _, address := range *NodeAddress {
				ipPortAlgo := fmt.Sprintf("-A %s %s:%d -s %s\n", proto, address, port.GetNodePort(), DefaultAlgo)
				svcString.WriteString(ipPortAlgo)

				endpoints := buildEndponts(address, proto, port.GetNodePort(), port.GetPort(), eps)
				svcString.WriteString(endpoints)

			}
		}
	}
	return svcString.String(), nil
}

func buildEndponts(VirtualIP, proto string, tgtPort, port int32, endpoints []*localnetv1.Endpoint) string {
	var strBuilder strings.Builder
	for _, endpoint := range endpoints {
		ip := endpoint.GetIPs().V4[0] //TODO:  This is what we call a Gambiarra :D Need to deal better with this thing.
		ep := fmt.Sprintf("-a %s %s:%d -r %s:%d -m -w %d\n", proto, VirtualIP, tgtPort, ip, port, DefaultWeight)
		strBuilder.WriteString(ep)
	}
	return strBuilder.String()
}
