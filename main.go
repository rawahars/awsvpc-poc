package main

import (
	"fmt"
	_ "plugin"
)

const (
	Add = true
	//Add = false
	networkNameENI = "task_br2"
	network2       = "taskbr02f17460b6c1"
	addressPrefix  = "10.0.0.228/24"
	ipAddr         = "10.0.0.228"

	adapterName = "Ethernet 6"

	containerID1    = "bd4fc4c9beb00070b86b78e7fe33cf9d340b6998f8e9338c1ad897b290bcf607" //"0a6940de65e0c11cf051cdd9d657e7dd89763ca5ca0b13551fffb7754a18119c"
	containerID2    = "328575d08fe2266f48a630f81ceabcb69b7a8c2f3383afecb3a46da262ee61f2" //"fd95ce38dbf4a35243edb8fc80d820ec3bfc0d3e249e89932aeb45808d13d891"
	endpointName    = "task-eni-ep2"
	natEndpointName = "nat-cid-12345"

	networkNameNAT = "nat"
	gateway        = "10.0.0.1"
)

func main() {
	if Add {
		CreateNetwork(networkNameENI, addressPrefix, gateway, adapterName)
		//CreateNetwork(networkNameNAT, addressPrefix, gateway)
		fmt.Println("Creating Endpoint now")
		CreateEndpoint(containerID1, true, endpointName, networkNameENI, gateway, ipAddr, false)
		CreateEndpoint(containerID2, false, endpointName, networkNameENI, gateway, ipAddr, false)
		//CreateEndpoint(containerID1, true, natEndpointName, networkNameNAT, gateway, ipAddr, true)
		fmt.Println("\nEndpoint Created and connected")
	} else {
		DeleteNetwork(networkNameENI)
		//DeleteNetwork(network2)
	}
}
