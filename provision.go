package upcloud

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"
	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"

	api_provision "github.com/wunderkraut/radi-api/operation/provision"
)

/**
 * Functionality for provisioning
 */

/**
 * HANDLER
 */

// UpCloud Provisioning Handler
type UpcloudProvisionHandler struct {
	BaseUpcloudServiceHandler
}

// Initialize and activate the Handler
func (provision *UpcloudProvisionHandler) Operations() api_operation.Operations {
	baseOperation := provision.BaseUpcloudServiceOperation()

	ops := api_operation.New_SimpleOperations()

	ops.Add(api_operation.Operation(&UpcloudProvisionUpOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudProvisionStopOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudProvisionDownOperation{BaseUpcloudServiceOperation: *baseOperation}))

	return ops.Operations()
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (provision *UpcloudProvisionHandler) Id() string {
	return "upcloud.provision"
}

/**
 * OPERATIONS
 */

// Provision up operation
type UpcloudProvisionUpOperation struct {
	BaseUpcloudServiceOperation
	api_provision.BaseProvisionUpOperation
}

// Return the string machinename/id of the Operation
func (up *UpcloudProvisionUpOperation) Id() string {
	return "upcloud.provision.up"
}

// Return a user readable string label for the Operation
func (up *UpcloudProvisionUpOperation) Label() string {
	return "Provision UpCloud servers"
}

// return a multiline string description for the Operation
func (up *UpcloudProvisionUpOperation) Description() string {
	return "Provision the UpCloud servers for this project."
}

// Run a validation check on the Operation
func (up *UpcloudProvisionUpOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// What settings/values does the Operation provide to an implemenentor
func (up *UpcloudProvisionUpOperation) Properties() api_property.Properties {
	return api_property.New_SimplePropertiesEmpty().Properties()
}

/**
 * Execute the Operation
 *
 * The following steps are followed for each server:
 *   1. create the server - then wait for it to be considered running
 *   2. create the firewall rules
 *   3. tag the server
 *
 * @TODO build properties properly from the child operations
 * @TODO This operation should operate in parrallel
 */
func (up *UpcloudProvisionUpOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	createOp := UpcloudServerCreateOperation{BaseUpcloudServiceOperation: up.BaseUpcloudServiceOperation}
	createProperties := createOp.Properties()

	service := up.ServiceWrapper()
	// settings := up.BuilderSettings()
	serverDefinitions := up.ServerDefinitions()

	// track which servers we actually create here
	createdServers := map[string]processedServer{}

	for _, id := range serverDefinitions.Order() {
		serverDefinition, _ := serverDefinitions.Get(id)
		createRequest := serverDefinition.CreateServerRequest()

		if requestProp, found := createProperties.Get(UPCLOUD_SERVER_CREATEREQUEST_PROPERTY); found {
			requestProp.Set(createRequest)
		}

		log.WithFields(log.Fields{"id": serverDefinition.Id()}).Info("Creating new server")

		createResult := createOp.Exec(createProperties)
		<-createResult.Finished()

		if !createResult.Success() {
			res.AddErrors(createResult.Errors())
			res.AddError(errors.New("Could not provision new UpCloud server: " + id))
			res.MarkFailed()
			continue
		} else {

			var createDetails upcloud.ServerDetails
			if detailsProp, found := createProperties.Get(UPCLOUD_SERVER_DETAILS_PROPERTY); found {
				createDetails = detailsProp.Get().(upcloud.ServerDetails)
			}

			uuid := createDetails.UUID

			createdServers[id] = processedServer{
				uuid:       uuid,
				definition: serverDefinition,
				details:    createDetails,
			}

			log.WithFields(log.Fields{"id": serverDefinition.Id(), "UUID": uuid, "state": createDetails.State}).Info("Created new server")
		}
	}

	firewallOp := UpcloudServerApplyFirewallRulesOperation{BaseUpcloudServiceOperation: up.BaseUpcloudServiceOperation}
	firewallProperties := firewallOp.Properties()

	// process tags and firewall rules
	for _, createdServer := range createdServers {
		uuid := createdServer.uuid
		serverDefinition := createdServer.definition

		// Before running anything, give the server a chance to get into the proper state
		log.WithFields(log.Fields{"id": serverDefinition.Id(), "UUID": uuid}).Info("Waiting for new server to start")
		if serverDetails, err := service.WaitForServerState(&upcloud_request.WaitForServerStateRequest{UUID: uuid, UndesiredState: "maintenance", Timeout: time.Minute * 2}); err != nil {
			if serverDetails != nil {
				uuid = serverDetails.UUID
			}
			res.AddError(err)
			res.AddError(errors.New("Server failed to start properly : " + uuid))
			res.MarkFailed()
		} else {
			log.WithFields(log.Fields{"state": serverDetails.State, "UUID": serverDetails.UUID}).Info("Server successfully created, now finalizing provisioning")

			serverDefinition := createdServer.definition
			firewallRules := serverDefinition.GetFirewallRules()

			if firewallProp, found := firewallProperties.Get(UPCLOUD_FIREWALL_RULES_PROPERTY); found {
				firewallProp.Set(firewallRules)
			}
			if uuidProp, found := firewallProperties.Get(UPCLOUD_SERVER_UUID_PROPERTY); found {
				uuidProp.Set(uuid)
			}

			firewallResult := firewallOp.Exec(firewallProperties)
			<-firewallResult.Finished()

			if !firewallResult.Success() {
				res.Merge(firewallResult)
				continue
			}

			// var serverDetails upcloud.ServerDetails
			// if detailsProp, found := createProperties.Get(UPCLOUD_SERVER_DETAILS_PROPERTY); found {
			// 	serverDetails = detailsProp.Get().(upcloud.ServerDetails)
			// }
		}
	}

	res.MarkFinished()

	return res.Result()
}

// Provision up operation
type UpcloudProvisionDownOperation struct {
	BaseUpcloudServiceOperation
	api_provision.BaseProvisionDownOperation
}

// Return the string machinename/id of the Operation
func (down *UpcloudProvisionDownOperation) Id() string {
	return "upcloud.provision.down"
}

// Return a user readable string label for the Operation
func (down *UpcloudProvisionDownOperation) Label() string {
	return "Remove UpCloud servers"
}

// return a multiline string description for the Operation
func (down *UpcloudProvisionDownOperation) Description() string {
	return "Remove the UpCloud servers for this project."
}

// Run a validation check on the Operation
func (down *UpcloudProvisionDownOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// What settings/values does the Operation provide to an implemenentor
func (down *UpcloudProvisionDownOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&UpcloudForceProperty{}))

	return props.Properties()
}

// Execute the Operation
//
// @TODO Add a way to remove the storage
// @TODO this operation could be optimized to work parrallel
func (down *UpcloudProvisionDownOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	downProperties := down.Properties()
	deleteOp := UpcloudServerDeleteOperation{BaseUpcloudServiceOperation: down.BaseUpcloudServiceOperation}
	deleteProperties := deleteOp.Properties()

	// service := down.ServiceWrapper()
	// settings := down.BuilderSettings()
	serverDefinitions := down.ServerDefinitions()

	// collect UUIDs of project servers
	uuids := []string{}
	for _, id := range serverDefinitions.Order() {
		serverDefinition, _ := serverDefinitions.Get(id)

		if serverDefinition.IsCreated() {
			uuid, _ := serverDefinition.UUID()
			log.WithFields(log.Fields{"id": id, "uuid": uuid}).Debug("Down: Server added to list")
			uuids = append(uuids, uuid)
		} else {
			log.WithFields(log.Fields{"id": id}).Info("Down: Server has not been created, so it will be skipped")
		}
	}

	if len(uuids) > 0 {

		if uuidsProp, found := deleteProperties.Get(UPCLOUD_SERVER_UUIDS_PROPERTY); found {
			log.WithFields(log.Fields{"uuids": uuids}).Info("DOWN: Using UUIDs")
			uuidsProp.Set(uuids)
		}
		if downForceProp, found := downProperties.Get(UPCLOUD_FORCE_PROPERTY); found {
			if deleteForceProp, found := deleteProperties.Get(UPCLOUD_FORCE_PROPERTY); found {
				if downForceProp.Get().(bool) {
					log.Info("DOWN: Forcing operation")
					deleteForceProp.Set(true)
				}
			}
		}

		log.WithFields(log.Fields{"uuids": uuids}).Info("Downing project servers")

		downResult := deleteOp.Exec(downProperties)
		<-downResult.Finished()

		res.Merge(downResult)

	} else {
		log.Info("No active servers found to take down.")
	}

	res.MarkFinished()

	return res.Result()
}

// Provision up operation
type UpcloudProvisionStopOperation struct {
	BaseUpcloudServiceOperation
	api_provision.BaseProvisionStopOperation
}

// Return the string machinename/id of the Operation
func (stop *UpcloudProvisionStopOperation) Id() string {
	return "upcloud.provision.stop"
}

// Return a user readable string label for the Operation
func (stop *UpcloudProvisionStopOperation) Label() string {
	return "Stop UpCloud servers"
}

// return a multiline string description for the Operation
func (stop *UpcloudProvisionStopOperation) Description() string {
	return "Stop the UpCloud servers for this project."
}

// Run a validation check on the Operation
func (stop *UpcloudProvisionStopOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// What settings/values does the Operation provide to an implemenentor
func (stop *UpcloudProvisionStopOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&UpcloudGlobalProperty{}))
	props.Add(api_property.Property(&UpcloudWaitProperty{}))
	props.Add(api_property.Property(&UpcloudServerUUIDProperty{}))

	return props.Properties()
}

// Execute the Operation
/**
 * @NOTE this is a first version.
 *
 * We will want to :
 *  1. retrieve servers by tag
 *  2. have a "remove-specific-uuid" option?
 */
func (stop *UpcloudProvisionStopOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	return res.Result()
}

// hold info about a server that we have processed
type processedServer struct {
	uuid       string
	definition ServerDefinition
	details    upcloud.ServerDetails
}
