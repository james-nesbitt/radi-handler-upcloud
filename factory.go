package upcloud

import (
	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"
	upcloud_client "github.com/Jalle19/upcloud-go-sdk/upcloud/client"
	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"
	upcloud_service "github.com/Jalle19/upcloud-go-sdk/upcloud/service"
)

const (
	CONFIG_KEY_UPCLOUD = "upcloud"
)

// A backend for pulling Config for UpCloud configuration
type UpcloudFactory interface {
	ServiceWrapper() *UpcloudServiceWrapper
	ServerDefinitions() ServerDefinitions
}

// Definition for a single UpCloud server
type ServerDefinition interface {
	Id() string
	UUID() (string, error)

	CreateServerRequest() upcloud_request.CreateServerRequest

	GetFirewallRules() upcloud.FirewallRules
	GetStorageDefinitions() StorageDefinitions

	GetServerDetails() (*upcloud.ServerDetails, error)
	GetServerState() (string, error)

	IsCreated() bool
	IsRunning() bool
}

// An ordered list of server definitions
type ServerDefinitions struct {
	defs  map[string]ServerDefinition
	order []string
}

// safe lazy initialzier
func (defs *ServerDefinitions) safe() {
	if defs.defs == nil {
		defs.defs = map[string]ServerDefinition{}
		defs.order = []string{}
	}
}

// Add a server def
func (defs *ServerDefinitions) Add(server ServerDefinition) {
	defs.safe()
	id := server.Id()
	if _, exists := defs.defs[id]; !exists {
		defs.order = append(defs.order, id)
	}
	defs.defs[id] = server
}

// Retrieve a server def by id
func (defs *ServerDefinitions) Get(id string) (ServerDefinition, bool) {
	defs.safe()
	def, exists := defs.defs[id]
	return def, exists
}

// return the ordered def keys
func (defs *ServerDefinitions) Order() []string {
	defs.safe()
	return defs.order
}

type StorageDefinition interface {
	Id() string
	BackupRule() upcloud.BackupRule
}

type StorageDefinitions struct {
	defs  map[string]StorageDefinition
	order []string
}

// safe lazy initialzier
func (defs *StorageDefinitions) safe() {
	if defs.defs == nil {
		defs.defs = map[string]StorageDefinition{}
		defs.order = []string{}
	}
}

// Add a storage def
func (defs *StorageDefinitions) Add(storage StorageDefinition) {
	defs.safe()
	id := storage.Id()
	if _, exists := defs.defs[id]; !exists {
		defs.order = append(defs.order, id)
	}
	defs.defs[id] = storage
}

// Retrieve a storage def by id
func (defs *StorageDefinitions) Get(id string) (StorageDefinition, bool) {
	defs.safe()
	def, exists := defs.defs[id]
	return def, exists
}

// return the ordered def keys
func (defs *StorageDefinitions) Order() []string {
	defs.safe()
	return defs.order
}

// Small factory used to create the ServiceWrapper from an UpCloud Client
type UpcloudServiceWrapperFactory struct {
	client upcloud_client.Client
}

// Constructor for UpcloudServiceWrapperFactory
func New_UpcloudServiceWrapperFactory(client upcloud_client.Client) *UpcloudServiceWrapperFactory {
	return &UpcloudServiceWrapperFactory{
		client: client,
	}
}

// Get an Upcloud service from these settings
func (serviceFactory UpcloudServiceWrapperFactory) Service() *upcloud_service.Service {
	return New_UpcloudServiceFromClient(serviceFactory.client)
}

// Get an Upcloud service from these settings
func (serviceFactory UpcloudServiceWrapperFactory) ServiceWrapper() *UpcloudServiceWrapper {
	service := serviceFactory.Service()
	return New_UpcloudServiceWrapper(*service)
}

// Constructor for upcloud Service from a client
func New_UpcloudServiceFromClient(client upcloud_client.Client) *upcloud_service.Service {
	service := upcloud_service.New(&client)
	return service
}
