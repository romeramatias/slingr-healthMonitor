package service

import "slingr/dto"

type HealthMonitor interface {
	Monitor(resource dto.Resource) (bool, error)
	Check() dto.ServerResponse
}
