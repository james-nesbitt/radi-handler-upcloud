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
	api_usage "github.com/wunderkraut/radi-api/usage"
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
 * only add internal operations.  For upcloud provision ops, we
 * tend to build the related operations directly in other handlers
 */

// UpCloud Provisioning Handler
type UpcloudServerHandler struct {
	BaseUpcloudServiceHandler
}

// Return a string identifier for the Handler (not functionally needed yet)
func (server *UpcloudServerHandler) Id() string {
	return "upcloud.server"
}

// Initialize and activate the Handler
func (server *UpcloudServerHandler) Operations() api_operation.Operations {
	baseOperation := server.BaseUpcloudServiceOperation()

	ops := api_operation.New_SimpleOperations()

	ops.Add(api_operation.Operation(&UpcloudServerCreateOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudServerStopOperation{BaseUpcloudServiceOperation: *baseOperation}))
	ops.Add(api_operation.Operation(&UpcloudServerDeleteOperation{BaseUpcloudServiceOperation: *baseOperation}))

	return ops.Operations()
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

// return a multiline string man page for the Operation
func (create *UpcloudServerCreateOperation) Help() string {
	return ""
}

// Run a validation check on the Operation
func (create *UpcloudServerCreateOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Is this operation an internal Operation
func (create *UpcloudServerCreateOperation) Usage() api_usage.Usage {
	return api_operation.Usage_Internal()
}

// What settings/values does the Operation provide to an implemenentor
func (create *UpcloudServerCreateOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&UpcloudServerCreateRequestProperty{}))
	props.Add(api_property.Property(&UpcloudServerDetailsProperty{}))

	return props.Properties()
}

// Execute the Operation
/**
 * @note this is a first version of the operation.  It does not implement
 *   the following checks/functionality:
 *     1. are the servies already provisioned?
 *     2. get the servers defintions from settings
 */
func (create *UpcloudServerCreateOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

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

		res.MarkSuccess()
	} else {
		res.AddError(errors.New("Unable to provision new server."))
		res.MarkFailed()
	}

	res.MarkFinished()

	return res.Result()
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

// return a multiline string man page for the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Help() string {
	return ""
}

// Run a validation check on the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Is this operation an internal Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Usage() api_usage.Usage {
	return api_operation.Usage_Internal()
}

// What settings/values does the Operation provide to an implemenentor
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&UpcloudFirewallRulesProperty{}))
	props.Add(api_property.Property(&UpcloudServerUUIDProperty{}))
	props.Add(api_property.Property(&UpcloudServerDetailsProperty{}))

	return props.Properties()
}

// Execute the Operation
func (applyFirewall *UpcloudServerApplyFirewallRulesOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

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
			res.AddError(err)
			res.MarkFailed()
		} else {
			log.WithFields(log.Fields{"position": ruleDetails.Position, "comment": ruleDetails.Comment, "uuid": uuid}).Info("Created server firewall rule")
			res.MarkSuccess()
		}
	}

	res.MarkFinished()

	return res.Result()
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

// return a multiline string man page for the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Help() string {
	return ""
}

// Run a validation check on the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Is this operation an internal Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Usage() api_usage.Usage {
	return api_operation.Usage_Internal()
}

// What settings/values does the Operation provide to an implemenentor
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&UpcloudStorageUUIDProperty{}))
	props.Add(api_property.Property(&UpcloudServerDetailsProperty{}))

	return props.Properties()
}

// Execute the Operation
func (applyBackup *UpcloudStorageApplyBackupRulesOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	// service := applyBackup.ServiceWrapper()
	// settings := applyBackup.BuilderSettings()

	// properties := applyBackup.Properties()

	return res.Result()
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

// return a multiline string man page for the Operation
func (delete *UpcloudServerDeleteOperation) Help() string {
	return ""
}

// Run a validation check on the Operation
func (delete *UpcloudServerDeleteOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Is this operation an internal Operation
func (delete *UpcloudServerDeleteOperation) Usage() api_usage.Usage {
	return api_operation.Usage_Internal()
}

// What settings/values does the Operation provide to an implemenentor
func (delete *UpcloudServerDeleteOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&UpcloudWaitProperty{}))
	props.Add(api_property.Property(&UpcloudForceProperty{}))
	props.Add(api_property.Property(&UpcloudServerUUIDSProperty{}))

	return props.Properties()
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
func (delete *UpcloudServerDeleteOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

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
				res.AddError(err)
				res.AddError(errors.New("Server not found, so cannot be deleted."))
				res.MarkFailed()
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
						res.AddError(err)
						res.AddError(errors.New("timeout waiting for server be removed."))
						res.MarkFailed()
					}
				} else {
					count++
					log.WithFields(log.Fields{"UUID": uuid}).Info("Removed UpCloud server")
				}
			} else {
				res.AddError(err)
				res.AddError(errors.New("Could not remove UpCloud server"))
				res.MarkFailed()
			}
		}

	} else {
		log.Info("No servers requested.  You should have passed a server UUID") // @TODO remove this when we are tagging servers
		res.MarkSuccess()
	}

	res.MarkFinished()

	return res.Result()
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

// return a multiline string man page for the Operation
func (stop *UpcloudServerStopOperation) Help() string {
	return ""
}

// Run a validation check on the Operation
func (stop *UpcloudServerStopOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Is this operation an internal Operation
func (stop *UpcloudServerStopOperation) Usage() api_usage.Usage {
	return api_operation.Usage_Internal()
}

// What settings/values does the Operation provide to an implemenentor
func (stop *UpcloudServerStopOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&UpcloudGlobalProperty{}))
	props.Add(api_property.Property(&UpcloudWaitProperty{}))
	props.Add(api_property.Property(&UpcloudServerUUIDSProperty{}))

	return props.Properties()
}

// Execute the Operation
/**
 * @NOTE this is a first version.
 */
func (stop *UpcloudServerStopOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

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
						res.AddError(err)
						res.AddError(errors.New("timeout waiting for server stop."))
						res.MarkFailed()
					}
				} else {
					log.WithFields(log.Fields{"UUID": uuid, "state": details.State, "progress": details.Progress}).Info("Stopped UpCloud server")
				}
			} else {
				res.AddError(err)
				res.AddError(errors.New("Could not stop UpCloud server"))
				res.MarkFailed()
			}
		}

	} else {
		log.Info("No servers requested.  You should have passed a server UUID") // @TODO remove this when we are tagging servers
	}

	return res.Result()
}
