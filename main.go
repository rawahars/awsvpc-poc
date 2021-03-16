package main

import (
	"fmt"
	_ "plugin"
)

const(
	Add = true
	//Add = false
	networkNameENI = "task_br"
	addressPrefix = "10.0.0.159/30"
	ipAddr = "10.0.0.159"

	containerID1 = "84a8b335171e246d6afd833324ad0fa1f0c6a549d00773b5292dc38cd9e1c161"
	containerID2 = ""
	endpointName = "task-eni-ep"


	networkNameNAT = "nat"
	gateway = "10.0.0.1"
)

func main(){
	if Add {
		CreateNetwork(networkNameENI, addressPrefix, gateway)
		fmt.Println("Creating Endpoint now")
		CreateEndpoint(containerID1, true, endpointName, networkNameENI, gateway, ipAddr)
		fmt.Println("\nEndpoint Created and connected")
	} else {
		DeleteNetwork(networkNameENI)
	}
}