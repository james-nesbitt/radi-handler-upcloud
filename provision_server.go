package upcloud

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"

	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"
	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"

	api_operation "github.com/wunderkraut/radi-api/operation"
)

/**
 * Here are a number of server provision related operations, none
 * of which are public, but all of which are used together in the
 * more public provision operations.
 */

/**
 * HANDLER
 *
 * @Note that this handler is not typically needed as it would
 * only add internal operations.  For upcloud provision ops, the
 * tend to build the related operations directly.
 */

// UpCloud Provisioning Handler
type UpcloudServerHandler struct {
	BaseUpcloudServiceHandler
}

// Initialize and activate the Handler
func (server *UpcloudServerHandler) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	baseOperation := server.BaseUpcloudServiceOperation()

	ops := api_operation.Operations{}

	ops.Add(api_operation.Operation(&UpcloudServerCreateOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudServerStopOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudServerDeleteOperation{BaseUpcloudServiceOperation: *baseOperation}))

	server.operations = &ops

	return api_operation.Result(result)
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (server *UpcloudServerHandler) Id() string {
	return "upcloud.server"
}

/**
 * OPERATIONS
 */

// Create a new server operation
type UpcloudServerCreateOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (create *UpcloudServerCreateOperation) Id() string {
	return "upcloud.server.create"
}

// Return a user readable string label for the Operation
func (create *UpcloudServerCreateOperation) Label() string {
	return "Create UpCloud server"
}

// return a multiline string description for the Operation
func (create *UpcloudServerCreateOperation) Description() string {
	return "Create an UpCloud server for this project."
}

// Run a validation check on the Operation
func (create *UpcloudServerCreateOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (create *UpcloudServerCreateOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (create *UpcloudServerCreateOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudServerCreateRequestProperty{}))
	props.Add(api_operation.Property(&UpcloudServerDetailsProperty{}))

	return props
}

// Execute the Operation
/**
 * @note this is a first version of the operation.  It does not implement
 *   the following checks/functionality:
 *     1. are the servies already provisioned?
 *     2. get the servers defintions from settings
 */
func (create *UpcloudServerCreateOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := create.ServiceWrapper()
	// settings := create.BuilderSettings()

	properties := create.Properties()

	request := upcloud_request.CreateServerRequest{}
	if requestProp, found := properties.Get(UPCLOUD_SERVER_CREATEREQUEST_PROPERTY); found {
		request = requestProp.Get().(upcloud_request.CreateServerRequest)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_CREATEREQUEST_PROPERTY, "prop": requestProp, "value": request}).Debug("Retrieved create server request")
	}

	log.WithFields(log.Fields{"request": request, "zone": request.Zone, "title": request.Title, "user": request.LoginUser}).Debug("Server: Using request to create a new server")
	serverDetails, err := service.CreateServer(&request)

	if err == nil {
		if detailsProp, found := properties.Get(UPCLOUD_SERVER_DETAILS_PROPERTY); found {
			detailsProp.Set(*serverDetails)
		}
		log.WithFields(log.Fields{"UUID": serverDetails.UUID}).Debug("server: Server created")

		result.MarkSuccess()
	} else {
		result.AddError(errors.New("Unable to provision new server."))
		result.MarkFailed()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

// Apply firewall rules to a running server
type UpcloudServerApplyFirewallRulesOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Id() string {
	return "upcloud.server.applyfirewallrules"
}

// Return a user readable string label for the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Label() string {
	return "Apply firewall rules"
}

// return a multiline string description for the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Description() string {
	return "Apply firewall rules to running UpCloud server."
}

// Run a validation check on the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudFirewallRulesProperty{}))
	props.Add(api_operation.Property(&UpcloudServerUUIDProperty{}))
	props.Add(api_operation.Property(&UpcloudServerDetailsProperty{}))

	return props
}

// Execute the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := applyFirewall.ServiceWrapper()
	// settings := applyFirewall.BuilderSettings()

	rules := upcloud.FirewallRules{}
	if rulesProp, found := props.Get(UPCLOUD_FIREWALL_RULES_PROPERTY); found {
		rules = rulesProp.Get().(upcloud.FirewallRules)
		log.WithFields(log.Fields{"key": UPCLOUD_FIREWALL_RULES_PROPERTY, "prop": rulesProp, "value": rules}).Debug("Retrieved firewall rules")
	}
	uuid := ""
	if uuidProp, found := props.Get(UPCLOUD_SERVER_UUID_PROPERTY); found {
		uuid = uuidProp.Get().(string)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_UUID_PROPERTY, "prop": uuidProp, "value": uuid}).Debug("Retrieved server UUID")
	}

	log.WithFields(log.Fields{"UUID": uuid, "#rules": len(rules.FirewallRules)}).Debug("Server: Applying firewall rules to server")

	for index, rule := range rules.FirewallRules {
		request := upcloud_request.CreateFirewallRuleRequest{
			FirewallRule: rule,
			ServerUUID:   uuid,
		}

		ruleDetails, err := service.CreateFirewallRule(&request)

		if err != nil {
			log.WithError(err).WithFields(log.Fields{"index": index, "position": rule.Position, "rule": rule, "rule-details": ruleDetails, "uuid": uuid}).Error("Failed to create server firewall rule")
			result.AddError(err)
			result.MarkFailed()
		} else {
			log.WithFields(log.Fields{"position": ruleDetails.Position, "comment": ruleDetails.Comment, "uuid": uuid}).Info("Created server firewall rule")
			result.MarkSuccess()
		}
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

// Apply firewall rules to a running server
type UpcloudStorageApplyBackupRulesOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Id() string {
	return "upcloud.storage.applybackuprules"
}

// Return a user readable string label for the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Label() string {
	return "Apply storage backup rules"
}

// return a multiline string description for the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Description() string {
	return "Apply storage backup rules"
}

// Run a validation check on the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudStorageUUIDProperty{}))
	props.Add(api_operation.Property(&UpcloudServerDetailsProperty{}))

	return props
}

// Execute the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	// service := applyBackup.ServiceWrapper()
	// settings := applyBackup.BuilderSettings()

	// properties := applyBackup.Properties()

	return api_operation.Result(result)
}

// Delete a server operation
type UpcloudServerDeleteOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (delete *UpcloudServerDeleteOperation) Id() string {
	return "upcloud.server.delete"
}

// Return a user readable string label for the Operation
func (delete *UpcloudServerDeleteOperation) Label() string {
	return "Remove UpCloud servers"
}

// return a multiline string description for the Operation
func (delete *UpcloudServerDeleteOperation) Description() string {
	return "Remove UpCloud servers for this project."
}

// Run a validation check on the Operation
func (delete *UpcloudServerDeleteOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (delete *UpcloudServerDeleteOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (delete *UpcloudServerDeleteOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudWaitProperty{}))
	props.Add(api_operation.Property(&UpcloudForceProperty{}))
	props.Add(api_operation.Property(&UpcloudServerUUIDSProperty{}))

	return props
}

// Execute the Operation
/**
 * @NOTE this is a first version.
 * @TODO this is a prime candidate for goroutines now that we have threaded options
 *
 * We will want to :
 *  1. retrieve servers by tag
 *  2. have a "remove-specific-uuid" option?
 */
func (delete *UpcloudServerDeleteOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := delete.ServiceWrapper()
	settings := delete.BuilderSettings()

	properties := delete.Properties()

	global := false
	if globalProp, found := properties.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("DELETE: Allowing global access")
	}
	wait := false
	if waitProp, found := properties.Get(UPCLOUD_WAIT_PROPERTY); found {
		wait = waitProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_WAIT_PROPERTY, "prop": waitProp, "value": wait}).Debug("DELETE: Wait for operation to complete")
	}
	force := false
	if waitProp, found := properties.Get(UPCLOUD_FORCE_PROPERTY); found {
		force = waitProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_FORCE_PROPERTY, "prop": waitProp, "value": force}).Debug("DELETE: force operation activated.")
	}
	uuidMatch := []string{}
	if uuidsProp, found := properties.Get(UPCLOUD_SERVER_UUIDS_PROPERTY); found {
		newUUIDs := uuidsProp.Get().([]string)
		uuidMatch = append(uuidMatch, newUUIDs...)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_UUIDS_PROPERTY, "prop": uuidsProp, "value": uuidMatch}).Debug("DELETE: Filter Server UUID")
	}

	if len(uuidMatch) > 0 {

		count := 0
		for _, uuid := range uuidMatch {
			if !(global || settings.ServerUUIDAllowed(uuid)) {
				log.WithFields(log.Fields{"uuid": uuid}).Error("Server UUID not a part of the project. Details will not be shown.")
				continue
			}

			details, err := service.GetServerDetails(&upcloud_request.GetServerDetailsRequest{UUID: uuid})

			if err != nil {
				result.AddError(err)
				result.AddError(errors.New("Server not found, so cannot be deleted."))
				result.MarkFailed()
				continue
			}

			if force && details.State == upcloud.ServerStateStarted {
				log.WithFields(log.Fields{"UUID": uuid, "state": details.State}).Warn("Stopping UpCloud server before deleting it.")
				_, err := service.StopServer(&upcloud_request.StopServerRequest{
					UUID:     details.UUID,
					StopType: upcloud_request.ServerStopTypeHard,
					Timeout:  time.Minute * 2,
				})
				if err != nil {
					log.WithFields(log.Fields{"UUID": uuid}).Warn("UpCloud server failed to stop before being deleted.")
					continue
				} else if waitDetails, err := service.WaitForServerState(&upcloud_request.WaitForServerStateRequest{UUID: uuid, DesiredState: upcloud.ServerStateStopped, Timeout: time.Minute * 2}); err != nil {
					log.WithFields(log.Fields{"UUID": uuid, "state": waitDetails.State}).Warn("UpCloud server failed to stop before being deleted.")
				}
			}

			request := upcloud_request.DeleteServerRequest{
				UUID: details.UUID,
			}
			err = service.DeleteServer(&request)

			if err == nil {
				if wait {
					waitRequest := upcloud_request.WaitForServerStateRequest{
						UUID:         uuid,
						DesiredState: "stopped",
						Timeout:      time.Duration(60) * time.Second,
					}
					details, err := service.WaitForServerState(&waitRequest)

					if err == nil {
						count++
						log.WithFields(log.Fields{"UUID": uuid, "state": details.State, "progress": details.Progress}).Info("Removed UpCloud server")
					} else {
						result.AddError(err)
						result.AddError(errors.New("timeout waiting for server be removed."))
						result.MarkFailed()
					}
				} else {
					count++
					log.WithFields(log.Fields{"UUID": uuid}).Info("Removed UpCloud server")
				}
			} else {
				result.AddError(err)
				result.AddError(errors.New("Could not remove UpCloud server"))
				result.MarkFailed()
			}
		}

	} else {
		log.Info("No servers requested.  You should have passed a server UUID") // @TODO remove this when we are tagging servers
		result.MarkSuccess()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

// Provision up operation
type UpcloudServerStopOperation struct {
	BaseUpcloudServiceOperation
}

// Return the string machinename/id of the Operation
func (stop *UpcloudServerStopOperation) Id() string {
	return "upcloud.server.stop"
}

// Return a user readable string label for the Operation
func (stop *UpcloudServerStopOperation) Label() string {
	return "Stop UpCloud server"
}

// return a multiline string description for the Operation
func (stop *UpcloudServerStopOperation) Description() string {
	return "Stop UpCloud servers."
}

// Run a validation check on the Operation
func (stop *UpcloudServerStopOperation) Validate() bool {
	return true
}

// Is this operation an internal Operation
func (stop *UpcloudServerStopOperation) Internal() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (stop *UpcloudServerStopOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&UpcloudGlobalProperty{}))
	props.Add(api_operation.Property(&UpcloudWaitProperty{}))
	props.Add(api_operation.Property(&UpcloudServerUUIDSProperty{}))

	return props
}

// Execute the Operation
/**
 * @NOTE this is a first version.
 */
func (stop *UpcloudServerStopOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := stop.ServiceWrapper()
	settings := stop.BuilderSettings()

	global := false
	if globalProp, found := props.Get(UPCLOUD_GLOBAL_PROPERTY); found {
		global = globalProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_GLOBAL_PROPERTY, "prop": globalProp, "value": global}).Debug("Allowing global access")
	}
	wait := false
	if waitProp, found := props.Get(UPCLOUD_WAIT_PROPERTY); found {
		wait = waitProp.Get().(bool)
		log.WithFields(log.Fields{"key": UPCLOUD_WAIT_PROPERTY, "prop": waitProp, "value": wait}).Debug("Wait for operation to complete")
	}
	uuidMatch := []string{}
	if uuidsProp, found := props.Get(UPCLOUD_SERVER_UUIDS_PROPERTY); found {
		newUUIDs := uuidsProp.Get().([]string)
		uuidMatch = append(uuidMatch, newUUIDs...)
		log.WithFields(log.Fields{"key": UPCLOUD_SERVER_UUIDS_PROPERTY, "prop": uuidsProp, "value": uuidMatch}).Debug("Filter: Server UUID")
	}

	if len(uuidMatch) > 0 {

		count := 0
		for _, uuid := range uuidMatch {
			if !(global || settings.ServerUUIDAllowed(uuid)) {
				log.WithFields(log.Fields{"uuid": uuid}).Error("Server UUID not a part of the project. Details will not be shown.")
				continue
			}

			request := upcloud_request.StopServerRequest{
				UUID: uuid,
			}

			log.WithFields(log.Fields{"uuid": uuid}).Info("Stopping server.")
			details, err := service.StopServer(&request)

			if err == nil {
				count++
				if wait {
					waitRequest := upcloud_request.WaitForServerStateRequest{
						UUID:           uuid,
						DesiredState:   "stopped",
						UndesiredState: "started",
						Timeout:        time.Duration(60) * time.Second,
					}
					details, err = service.WaitForServerState(&waitRequest)

					if err == nil {
						log.WithFields(log.Fields{"UUID": uuid, "state": details.State, "progress": details.Progress}).Info("Stopped UpCloud server")
					} else {
						result.AddError(err)
						result.AddError(errors.New("timeout waiting for server stop."))
						result.MarkFailed()
					}
				} else {
					log.WithFields(log.Fields{"UUID": uuid, "state": details.State, "progress": details.Progress}).Info("Stopped UpCloud server")
				}
			} else {
				result.AddError(err)
				result.AddError(errors.New("Could not stop UpCloud server"))
				result.MarkFailed()
			}
		}

	} else {
		log.Info("No servers requested.  You should have passed a server UUID") // @TODO remove this when we are tagging servers
	}

	return api_operation.Result(result)
}
