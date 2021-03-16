package main

import (
	"encoding/json"
	"fmt"
	"github.com/Microsoft/hcsshim"
	"net"
	"strings"
)

type hnsRoutePolicy struct {
	hcsshim.Policy
	DestinationPrefix string `json:"DestinationPrefix,omitempty"`
	NeedEncap         bool   `json:"NeedEncap,omitempty"`
}

func CreateNetwork(networkName string, addressPrefix string, gateway string){
	hnsNetwork, err := hcsshim.GetHNSNetworkByName(networkName)

	if err == nil {
		fmt.Println("Network present")
	} else {
		ipAddr, _ := GetIPAddressFromString(addressPrefix)
		// Initialize the HNS network.
		hnsNetwork = &hcsshim.HNSNetwork{
			Name:               networkName,
			Type:               "l2bridge",
			NetworkAdapterName: "",

			Subnets: []hcsshim.Subnet{
				{
					AddressPrefix:  GetSubnetPrefix(ipAddr).String(),
					GatewayAddress: gateway,
				},
			},
		}

		buf, err := json.Marshal(hnsNetwork)
		if err != nil {
			fmt.Println("Error Encountered while marshalling")
			return
		}
		hnsRequest := string(buf)

		hnsResponse, err := hcsshim.HNSNetworkRequest("POST", "", hnsRequest)
		if err != nil {
			fmt.Println("Error while network creation")
			return
		}
		fmt.Println("HNS Network is created")
		fmt.Printf(hnsResponse.Id)
	}
}

func DeleteNetwork(networkName string){
	// Find the HNS network ID.
	hnsNetwork, err := hcsshim.GetHNSNetworkByName(networkName)
	if err != nil {
		fmt.Printf("Error while network detection %s", err)
		return
	}

	// Delete the HNS network.
	fmt.Printf("Deleting HNS network name: %s ID: %s.", networkName, hnsNetwork.Id)
	_, err = hcsshim.HNSNetworkRequest("DELETE", hnsNetwork.Id, "")
	if err != nil {
		fmt.Printf("Failed to delete HNS network: %v.", err)
	}
}

func CreateEndpoint(containerID string, isInfraContainer bool, endpointName string, networkName string, gateway string,
	ipAddr string){
	hnsEndpoint, err := hcsshim.GetHNSEndpointByName(endpointName)
	if err == nil {
		fmt.Printf("Found Existing Endpoint: %s", endpointName)
		if !isInfraContainer{
			err := attachEndpoint(hnsEndpoint, containerID)
			if err != nil{
				fmt.Printf("Error while attaching endpoint")
				return
			}
		}
		return
	}
	fmt.Printf("Creating new endpoint with name %s", endpointName)

	DNSServers := [] string{"10.0.0.2"}
	// Initialize the HNS endpoint.
	hnsEndpoint = &hcsshim.HNSEndpoint{
		Name:               endpointName,
		VirtualNetworkName: networkName,
		DNSServerList:      strings.Join(DNSServers, ","),
		GatewayAddress:     gateway,
		IPAddress:          net.ParseIP(ipAddr),
	}

	// Attach policies
	err = addEndpointPolicy(
		hnsEndpoint,
		hnsRoutePolicy{
			Policy:            hcsshim.Policy{Type: hcsshim.Route},
			DestinationPrefix: ipAddr + "/32",
			NeedEncap:         true,
		})
	if err != nil {
		fmt.Printf("Failed to add endpoint route policy: %v.", err)
		return
	}

	err = addEndpointPolicy(
		hnsEndpoint,
		hnsRoutePolicy{
			Policy: hcsshim.Policy{Type: hcsshim.OutboundNat},
		})
	if err != nil {
		fmt.Printf("Failed to add OutboundNat policy: %v.", err)
		return
	}

	//////////////////
	// Encode the endpoint request.
	buf, err := json.Marshal(hnsEndpoint)
	if err != nil {
		return
	}
	hnsRequest := string(buf)

	// Create the HNS endpoint.
	fmt.Printf("Creating HNS endpoint: %+v", hnsRequest)
	hnsResponse, err := hcsshim.HNSEndpointRequest("POST", "", hnsRequest)
	if err != nil {
		fmt.Printf("Failed to create HNS endpoint: %v.", err)
		return
	}

	fmt.Printf("Received HNS endpoint response: %+v.", hnsResponse)

	// Attach the HNS endpoint to the container's network namespace.
	err = attachEndpoint(hnsResponse, containerID)
	if err != nil {
		// Cleanup the failed endpoint.
		fmt.Printf("Deleting the failed HNS endpoint %s.", hnsResponse.Id)
		_, delErr := hcsshim.HNSEndpointRequest("DELETE", hnsResponse.Id, "")
		if delErr != nil {
			fmt.Printf("Failed to delete HNS endpoint: %v.", delErr)
		}

		return
	}

	// Return network interface MAC address.
	mac, _ := net.ParseMAC(hnsResponse.MacAddress)
	fmt.Printf("MAC Address is %v", mac)

	return
}

func DeleteEndpoint(){

}

func attachEndpoint(ep *hcsshim.HNSEndpoint, containerID string) error {
	fmt.Printf("Attaching HNS endpoint %s to container %s.", ep.Id, containerID)
	err := hcsshim.HotAttachEndpoint(containerID, ep.Id)
	if err != nil {
		// Attach can fail if the container is no longer running and/or its network namespace
		// has been cleaned up.
		fmt.Printf("Failed to attach HNS endpoint %s: %v.", ep.Id, err)
	}

	return err
}

func addEndpointPolicy(ep *hcsshim.HNSEndpoint, policy interface{}) error {
	buf, err := json.Marshal(policy)
	if err != nil {
		fmt.Printf("Failed to encode policy: %v.", err)
		return err
	}

	ep.Policies = append(ep.Policies, buf)

	return nil
}