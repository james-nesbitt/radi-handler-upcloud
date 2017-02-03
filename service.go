package upcloud

import (
	upcloud_service "github.com/Jalle19/upcloud-go-sdk/upcloud/service"
)

/**
 * UpCloud SDK service wrapper
 */

/**
 * A Wrapper for the UpCloud Service so that we can add and streamline stuff
 */

// Constructor for UpcloudServiceWrapper
func New_UpcloudServiceWrapper(service upcloud_service.Service) *UpcloudServiceWrapper {
	return &UpcloudServiceWrapper{
		Service: service,
	}
}

// Wrapper for the upcloud service, so that we can limit operations
type UpcloudServiceWrapper struct {
	upcloud_service.Service
}
