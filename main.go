package main

import (
	"fmt"
	"slingr/dto"
	"slingr/service"
)

func main() {
	// Create a single instance of Health Monitor
	var healthMonitor = service.NewHealthMonitorImpl()

	// Try to insert resources, if healthMonitor doesn't have resources there is nothing to check
	// If type, name or handle are empty, the resource is not added
	resourceOne := dto.Resource{
		Type:     "serviceUrl",
		Name:     "graphql",
		Handle:   "uri",
		Critical: true,
	}
	_, errMonitorOne := healthMonitor.Monitor(resourceOne)
	if errMonitorOne != nil {
		fmt.Printf(errMonitorOne.Error())
	}

	resourceTwo := dto.Resource{
		Type:     "serviceUrl",
		Name:     "resource2",
		Handle:   "uri2",
		Critical: true,
	}

	_, errMonitorTwo := healthMonitor.Monitor(resourceTwo)
	if errMonitorTwo != nil {
		fmt.Printf(errMonitorTwo.Error())
	}

	// Check function
	serverResponse := healthMonitor.Check()
	if serverResponse.Status == 200 {
		fmt.Printf("[SERVER RESPONSE OK][STATUS: %v][MESSAGE: %v]\n", serverResponse.Status, serverResponse.Message)
		for _, s := range serverResponse.ServiceResponses {
			fmt.Printf("[SERVER RESPONSE OK][RESOURCE: %v][STATUS: %v]\n", s.Resource, s.Status)
		}
	} else {
		fmt.Printf("[SERVER ERROR][STATUS: %v][MESSAGE: %v]\n", serverResponse.Status, serverResponse.Message)
	}

}
