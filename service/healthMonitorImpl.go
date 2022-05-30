package service

import (
	"errors"
	"fmt"
	"slingr/dto"
	"slingr/utils"
	"sync"
	"time"
)

var once sync.Once
var healthMonitor HealthMonitor

type HealthMonitorImpl struct {
	monitors          map[string]map[string]string
	criticalResources []string
}

const (
	FuncServiceUrl            = "serviceUrl"
	FuncRedisClient           = "redisClient"
	FuncElasticsearchClient   = "elasticsearchClient"
	FuncPostgresPromiseClient = "postgresPromiseClient"
	// Timeout Can change the value to get a timeout, the services have 3 seconds response delay
	// if Timeout < time.Duration(3) throws timeout error
	Timeout = time.Duration(5)
)

func NewHealthMonitorImpl() HealthMonitor {
	if healthMonitor == nil {
		once.Do(
			func() {
				fmt.Println("[HealthMonitorImpl][Method: NewHealthMonitorImpl][Message: Creating instance...]")
				healthMonitor = &HealthMonitorImpl{monitors: make(map[string]map[string]string)}
			})
	} else {
		fmt.Println("[HealthMonitorImpl][Method: NewHealthMonitorImpl][Message: HealthMonitor already created]")
	}

	return healthMonitor
}

func (h *HealthMonitorImpl) Monitor(resource dto.Resource) (bool, error) {
	fmt.Printf("[HealthMonitorImpl][Method: Monitor][INIT][Params: %v]\n", resource)
	defer fmt.Printf("[HealthMonitorImpl][Method: Monitor][END]\n")

	if len(resource.Type) == 0 || len(resource.Name) == 0 || len(resource.Handle) == 0 {
		fmt.Printf("[HealthMonitorImpl][Method: Monitor][Failed to add resource]\n")
		return false, errors.New("error: resource with empty values")
	}

	resourceByType := h.monitors[resource.Type]
	if resourceByType == nil {
		resourceByType = make(map[string]string)
	}
	resourceByType[resource.Name] = resource.Handle
	h.monitors[resource.Type] = resourceByType

	if resource.Critical {
		h.criticalResources = append(h.criticalResources, resource.Name)
	}

	return true, nil
}

func (h *HealthMonitorImpl) Check() dto.ServerResponse {
	fmt.Printf("[HealthMonitorImpl][Method: Check][INIT]\n")
	defer fmt.Printf("[HealthMonitorImpl][Method: Check][END]\n")

	channelsArray, errMonitorsToCheck := h.monitorsToCheck()
	if errMonitorsToCheck != nil {
		return dto.ServerResponse{
			Status:  500,
			Message: "Nothing to check",
		}
	}
	fmt.Printf("[HealthMonitorImpl][Method: Check][Monitors to check: %v]\n", len(channelsArray))

	serviceResponses, errServiceResponses := h.getServiceResponses(channelsArray)
	if errServiceResponses != nil {
		if errServiceResponses.Error() == "timeout" {
			return dto.ServerResponse{
				Status:  503,
				Message: "Timed out while checking resources",
			}
		}

		return dto.ServerResponse{
			Status:  500,
			Message: "Generic error",
		}
	}
	fmt.Printf("[HealthMonitorImpl][Method: Check][Services responses: %v]\n", len(serviceResponses))

	serverResponse := dto.ServerResponse{
		Status:           200,
		ServiceResponses: serviceResponses,
		Message:          "Ok",
	}

	var failed []string
	for _, result := range serviceResponses {
		if result.Status != "ok" {
			if utils.ExistsInArray(h.criticalResources, result.Resource) {
				serverResponse.Status = 503
				serverResponse.Message = "Fail in a critical resource service"
			}
			failed = append(failed, result.Resource)
		}
	}

	if len(failed) > 0 {
		serverResponse.Failed = failed
	}

	return serverResponse
}

func (h *HealthMonitorImpl) monitorsToCheck() ([]<-chan dto.ServiceResponse, error) {
	fmt.Printf("[HealthMonitorImpl][Method: monitorsToCheck][INIT]\n")
	defer fmt.Printf("[HealthMonitorImpl][Method: monitorsToCheck][END]\n")

	var channelsArray []<-chan dto.ServiceResponse
	var function func(string, string) <-chan dto.ServiceResponse
	for key, value := range h.monitors {
		switch key {
		case FuncServiceUrl:
			function = h.checkServiceUrl
		case FuncRedisClient:
			function = h.checkRedisClient
		case FuncElasticsearchClient:
			function = h.checkElasticClient
		case FuncPostgresPromiseClient:
			function = h.checkPostgresPromiseClient
		}

		for k, v := range value {
			channelsArray = append(channelsArray, function(k, v))
		}
	}

	if channelsArray == nil || len(channelsArray) < 1 {
		return channelsArray, errors.New("nothing to check")
	}

	return channelsArray, nil
}

func (h *HealthMonitorImpl) getServiceResponses(channelsArray []<-chan dto.ServiceResponse) ([]dto.ServiceResponse, error) {
	fmt.Printf("[HealthMonitorImpl][Method: getServiceResponses][INIT]\n")
	defer fmt.Printf("[HealthMonitorImpl][Method: getServiceResponses][END]\n")

	var timeout = time.After(Timeout * time.Second)
	var serviceResponses []dto.ServiceResponse
	for _, c := range channelsArray {
		select {
		case res := <-c:
			serviceResponses = append(serviceResponses, res)
		case <-timeout:
			return serviceResponses, errors.New("timeout")
		}
	}

	return serviceResponses, nil
}

// HealthMonitor other services/clients functions
func (h *HealthMonitorImpl) checkServiceUrl(resourceName string, url string) <-chan dto.ServiceResponse {
	response := make(chan dto.ServiceResponse)

	go func() {
		// Can change the time value or the status value to get an error
		defer close(response)
		time.Sleep(time.Second * 3)
		response <- dto.ServiceResponse{
			Resource: resourceName,
			Status:   "ok",
		}
	}()

	return response
}

func (h *HealthMonitorImpl) checkRedisClient(resourceName string, redisClient string) <-chan dto.ServiceResponse {
	response := make(chan dto.ServiceResponse)

	go func() {
		defer close(response)
		time.Sleep(time.Second * 3)
		response <- dto.ServiceResponse{
			Resource: resourceName,
			Status:   "ok",
		}
	}()

	return response
}

func (h *HealthMonitorImpl) checkPostgresPromiseClient(resourceName string, postgresClient string) <-chan dto.ServiceResponse {
	response := make(chan dto.ServiceResponse)

	go func() {
		defer close(response)
		time.Sleep(time.Second * 3)
		response <- dto.ServiceResponse{
			Resource: resourceName,
			Status:   "ok",
		}
	}()

	return response
}

func (h *HealthMonitorImpl) checkElasticClient(resourceName string, elasticClient string) <-chan dto.ServiceResponse {
	response := make(chan dto.ServiceResponse)

	go func() {
		defer close(response)
		time.Sleep(time.Second * 3)
		response <- dto.ServiceResponse{
			Resource: resourceName,
			Status:   "ok",
		}
	}()

	return response
}
