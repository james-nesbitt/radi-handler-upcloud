package upcloud

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"

	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Monitor handler for Upcloud operations
 */
type UpcloudMonitorHandler struct {
	BaseUpcloudServiceHandler
}

// Initialize and activate the Handler
func (monitor *UpcloudMonitorHandler) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	baseOperation := monitor.BaseUpcloudServiceOperation()

	ops := api_operation.Operations{}
	ops.Add(api_operation.Operation(&UpcloudMonitorListZonesOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudMonitorListPlansOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudMonitorListServersOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudMonitorServerDetailsOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudMonitorListStoragesOperation{BaseUpcloudServiceOperation: *baseOperation}))
	monitor.operations = &ops

	return api_operation.Result(result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (monitor *UpcloudMonitorHandler) Id() string {
	return "upcloud.monitor"
}

/**
 * Monitor operations for UpCloud
 */

/**
 * List UpCloud Zones
 */
type UpcloudMonitorListZonesOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (listZones *UpcloudMonitorListZonesOperation) Id() string {
	return "upcloud.monitor.list.zones"
}

// Return a user readable string label for the Operation
func (listZones *UpcloudMonitorListZonesOperation) Label() string {
	return "List UpCloud Zones"
}

// return a multiline string description for the Operation
func (listZones *UpcloudMonitorListZonesOperation) Description() string {
	return "List information about the UpCloud zones."
}

// Is this operation meant to be used only inside the API
func (listZones *UpcloudMonitorListZonesOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (listZones *UpcloudMonitorListZonesOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (listZones *UpcloudMonitorListZonesOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
	props.Add(api_operation.Property(&UpcloudZoneIdProperty{}))

	return props
}

// Execute the Operation
func (listZones *UpcloudMonitorListZonesOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := listZones.ServiceWrapper()
	// settings := listZones.BuilderSettings()

	global := false
	if globalProp, found := props.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("Filter: global")
	}
	idMatch := []string{}
	if idProp, found := props.Get(UPCLOUD_ZONE_ID_PROPERTY); found {
		ids := idProp.Get().([]string)
		idMatch = append(idMatch, ids...)
		log.WithFields(log.Fields{"key": UPCLOUD_ZONE_ID_PROPERTY, "prop": idProp, "value": idMatch}).Debug("Filter: zone id")
	}

	zones, err := service.GetZones()
	if err == nil {
		for index, zone := range zones.Zones {
			filterOut := false

			// if some server filters were passed, filter out anything not in the passed list
			if len(idMatch) > 0 {
				found := false
				for _, id := range idMatch {
					if id == zone.Id {
						found = true
					}
				}
				filterOut = !found
			}

			if !filterOut {
				log.WithFields(log.Fields{"index": index, "id": zone.Id, "description": zone.Description}).Info("UpCloud zone")
			}
		}

		result.MarkSuccess()
	} else {
		result.AddError(err)
		result.AddError(errors.New("Could not retrieve UpCloud zones information."))
		result.MarkFailed()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

/**
 * Monitor operations for UpCloud
 */
type UpcloudMonitorListServersOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (listServers *UpcloudMonitorListServersOperation) Id() string {
	return "upcloud.monitor.list.servers"
}

// Return a user readable string label for the Operation
func (listServers *UpcloudMonitorListServersOperation) Label() string {
	return "UpCloud server list"
}

// return a multiline string description for the Operation
func (listServers *UpcloudMonitorListServersOperation) Description() string {
	return "List UpCloud servers used in the project"
}

// Is this operation meant to be used only inside the API
func (listServers *UpcloudMonitorListServersOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (listServers *UpcloudMonitorListServersOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (listServers *UpcloudMonitorListServersOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
	props.Add(api_operation.Property(&UpcloudServerUUIDSProperty{}))

	return props
}

// Execute the Operation
func (listServers *UpcloudMonitorListServersOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := listServers.ServiceWrapper()
	settings := listServers.BuilderSettings()
	serverDefinitions := listServers.ServerDefinitions()

	projectUUIDs := []string{}
	for _, id := range serverDefinitions.Order() {
		serverDefinition, _ := serverDefinitions.Get(id)

		if serverDefinition.IsCreated() {
			uuid, _ := serverDefinition.UUID()
			log.WithFields(log.Fields{"id": id, "uuid": uuid}).Debug("Monitor: Server added to list")
			projectUUIDs = append(projectUUIDs, uuid)
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Monitor: Server has not been created, so it will be skipped")
		}
	}

	global := false
	properties := listServers.Properties()
	if globalProp, found := properties.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("GLOBAL")
	}
	uuidMatch := []string{}
	if uuidProp, found := properties.Get(UPCLOUD_SERVER_UUIDS_PROPERTY); found {
		for _, newUUID := range uuidProp.Get().([]string) {
			if global {
				uuidMatch = append(uuidMatch, newUUID)
			} else {
				for _, projectUUID := range projectUUIDs {
					if projectUUID == newUUID {
						uuidMatch = append(uuidMatch, newUUID)
						break
					}
				}
			}
		}
	}

	servers, err := service.GetServers()
	if err == nil {
		serverList := servers.Servers
		if len(serverList) > 0 {
			for index, server := range serverList {
				filterOut := false

				// filter out servers that are no a part of the current project
				if !global {
					filterOut = !settings.ServerUUIDAllowed(server.UUID)
				}

				// if some server filters were passed, filter out anything not in the passed list
				if len(uuidMatch) > 0 {
					found := false
					for _, uuid := range uuidMatch {
						if uuid == server.UUID {
							found = true
						}
					}
					filterOut = !found
				}

				if !filterOut {
					log.WithFields(log.Fields{"index": index, "uuid": server.UUID, "title": server.Title, "plan": server.Plan, "zone": server.Zone, "state": server.State, "progress": server.Progress, "tags": server.Tags}).Info("Server")
				}
			}
		} else {
			log.WithFields(log.Fields{"Filter UUIDs": uuidMatch}).Info("No servers found")
		}

		result.MarkSuccess()
	} else {
		result.AddError(err)
		result.AddError(errors.New("Could not retrieve upcloud server list."))
		result.MarkFailed()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

/**
 * Monitor operations for UpCloud
 */
type UpcloudMonitorServerDetailsOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (serverDetail *UpcloudMonitorServerDetailsOperation) Id() string {
	return "upcloud.monitor.server.details"
}

// Return a user readable string label for the Operation
func (serverDetail *UpcloudMonitorServerDetailsOperation) Label() string {
	return "UpCloud server details"
}

// return a multiline string description for the Operation
func (serverDetail *UpcloudMonitorServerDetailsOperation) Description() string {
	return "UpCloud server details"
}

// Is this operation meant to be used only inside the API
func (serverDetail *UpcloudMonitorServerDetailsOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (serverDetail *UpcloudMonitorServerDetailsOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (serverDetail *UpcloudMonitorServerDetailsOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
	props.Add(api_operation.Property(&UpcloudServerUUIDSProperty{}))

	return props
}

// Execute the Operation
func (serverDetail *UpcloudMonitorServerDetailsOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := serverDetail.ServiceWrapper()
	settings := serverDetail.BuilderSettings()
	serverDefinitions := serverDetail.ServerDefinitions()

	projectUUIDs := []string{}
	for _, id := range serverDefinitions.Order() {
		serverDefinition, _ := serverDefinitions.Get(id)

		if serverDefinition.IsCreated() {
			uuid, _ := serverDefinition.UUID()
			log.WithFields(log.Fields{"id": id, "uuid": uuid}).Debug("Monitor: Server added to list")
			projectUUIDs = append(projectUUIDs, uuid)
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Monitor: Server has not been created, so it will be skipped")
		}
	}

	global := false
	if globalProp, found := props.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("GLOBAL")
	}
	uuidMatch := []string{}
	if uuidProp, found := props.Get(UPCLOUD_SERVER_UUIDS_PROPERTY); found {
		for _, newUUID := range uuidProp.Get().([]string) {
			if global {
				uuidMatch = append(uuidMatch, newUUID)
			} else {
				for _, projectUUID := range projectUUIDs {
					if projectUUID == newUUID {
						uuidMatch = append(uuidMatch, newUUID)
						break
					}
				}
			}
		}
	}
	if len(uuidMatch) == 0 {
		uuidMatch = projectUUIDs
		log.WithFields(log.Fields{"uuids": uuidMatch}).Debug("Filter: Server UUIDs")
	} else {
		log.WithFields(log.Fields{"uuids": uuidMatch}).Debug("Filter: All Server UUIDs")
	}

	if len(uuidMatch) > 0 {

		count := 0
		for _, uuid := range uuidMatch {
			if !(global || settings.ServerUUIDAllowed(uuid)) {
				log.WithFields(log.Fields{"uuid": uuid}).Error("Server UUID not a part of the project. Details will not be shown.")
				continue
			}

			request := upcloud_request.GetServerDetailsRequest{UUID: uuid}

			if details, err := service.GetServerDetails(&request); err == nil {
				count++
				log.WithFields(log.Fields{"index": count, "UUID": uuid, "tags": details.Tags, "details": details}).Info("Server Details")
			} else {
				log.WithError(err).WithFields(log.Fields{"UUID": uuid}).Error("Could not fetch server details.")
				result.AddError(err)
			}
		}

		if count == 0 {
			result.AddError(errors.New("No servers were matched."))
		}

		result.MarkSuccess()
	} else {
		result.AddError(errors.New("No servers uuids were passed to monitor server details operation, so no details can be shown."))
		result.MarkFailed()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

/**
 * Monitor operations for UpCloud
 */
type UpcloudMonitorListPlansOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (listPlans *UpcloudMonitorListPlansOperation) Id() string {
	return "upcloud.monitor.list.plans"
}

// Return a user readable string label for the Operation
func (listPlans *UpcloudMonitorListPlansOperation) Label() string {
	return "UpCloud server plans"
}

// return a multiline string description for the Operation
func (listPlans *UpcloudMonitorListPlansOperation) Description() string {
	return "List UpCloud plans avaialble"
}

// Is this operation meant to be used only inside the API
func (listPlans *UpcloudMonitorListPlansOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (listPlans *UpcloudMonitorListPlansOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (listPlans *UpcloudMonitorListPlansOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	// props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
	// props.Add(api_operation.Property(&UpcloudServerUUIDProperty{}))

	return props
}

// Execute the Operation
func (listPlans *UpcloudMonitorListPlansOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := listPlans.ServiceWrapper()
	//settings := listPlans.BuilderSettings()

	plans, err := service.GetPlans()
	if err == nil {
		if len(plans.Plans) > 0 {
			for index, plan := range plans.Plans {
				log.WithFields(log.Fields{"index": index, "name": plan.Name, "plan": plan}).Info("Plan")
			}
		} else {
			log.Info("No plans are available")
		}

	} else {
		result.AddError(err)
		result.MarkFailed()
	}

	return api_operation.Result(result)
}

/**
 * Monitor operations for UpCloud
 */
type UpcloudMonitorListStoragesOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (listStorages *UpcloudMonitorListStoragesOperation) Id() string {
	return "upcloud.monitor.list.storages"
}

// Return a user readable string label for the Operation
func (listStorages *UpcloudMonitorListStoragesOperation) Label() string {
	return "UpCloud storage list"
}

// return a multiline string description for the Operation
func (listStorages *UpcloudMonitorListStoragesOperation) Description() string {
	return "List UpCloud storages"
}

// Is this operation meant to be used only inside the API
func (listStorages *UpcloudMonitorListStoragesOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (listStorages *UpcloudMonitorListStoragesOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (listStorages *UpcloudMonitorListStoragesOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
	props.Add(api_operation.Property(&UpcloudStorageUUIDProperty{}))

	return props
}

// Execute the Operation
func (listStorages *UpcloudMonitorListStoragesOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := listStorages.ServiceWrapper()
	settings := listStorages.BuilderSettings()

	global := false
	if globalProp, found := props.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("Filter: Global")
	}
	uuidMatch := []string{}
	if uuidProp, found := props.Get(UPCLOUD_STORAGE_UUID_PROPERTY); found {
		newUUIDs := uuidProp.Get().([]string)
		uuidMatch = append(uuidMatch, newUUIDs...)
		log.WithFields(log.Fields{"key": UPCLOUD_STORAGE_UUID_PROPERTY, "prop": uuidMatch, "value": uuidMatch}).Debug("Filter: Storage UUID")
	}

	request := upcloud_request.GetStoragesRequest{}

	storages, err := service.GetStorages(&request)
	if err == nil {
		storageList := storages.Storages
		if len(storageList) > 0 {
			for index, storage := range storageList {
				filterOut := false

				// filter out servers that are no a part of the current project
				if !global {
					filterOut = !settings.StorageUUIDAllowed(storage.UUID)
				}

				// if some storage filters were passed, filter out anything not in the passed list
				if len(uuidMatch) > 0 {
					found := false
					for _, uuid := range uuidMatch {
						if uuid == storage.UUID {
							found = true
							break
						}
					}
					filterOut = !found
				}

				if !filterOut {
					log.WithFields(log.Fields{"index": index, "uuid": storage.UUID, "title": storage.Title, "type": storage.Type, "plan": storage.PartOfPlan, "zone": storage.Zone, "size": storage.Size}).Info("Storage")
				}
			}
		} else {
			log.WithFields(log.Fields{"acces": request.Access, "type": request.Type, "favoriate": request.Favorite}).Info("No storages found")
		}

		result.MarkSuccess()
	} else {
		result.AddError(errors.New("Could not retrieve upcloud storage list."))
		result.MarkFailed()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
