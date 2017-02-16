package upcloud

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"

	api_security "github.com/wunderkraut/radi-api/operation/security"
)

/**
 * Some security implementations for upcloud
 */

/**
 * Security handler for Upcloud operations
 */
type UpcloudSecurityHandler struct {
	BaseUpcloudServiceHandler
}

// Initialize and activate the Handler
func (security *UpcloudSecurityHandler) Operations() api_operation.Operations {
	baseOperation := security.BaseUpcloudServiceOperation()

	ops := api_operation.New_SimpleOperations()
	ops.Add(api_operation.Operation(&UpcloudSecurityUserOperation{BaseUpcloudServiceOperation: *baseOperation}))

	return ops.Operations()
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (security *UpcloudSecurityHandler) Id() string {
	return "upcloud.security"
}

// A security user information operation
type UpcloudSecurityUserOperation struct {
	BaseUpcloudServiceOperation
	api_security.BaseSecurityUserOperation
}

// Return the string machinename/id of the Operation
func (securityUser *UpcloudSecurityUserOperation) Id() string {
	return "upcloud.security.account"
}

// Return a user readable string label for the Operation
func (securityUser *UpcloudSecurityUserOperation) Label() string {
	return "Show UpCloud Account information"
}

// return a multiline string description for the Operation
func (securityUser *UpcloudSecurityUserOperation) Description() string {
	return "Show information about the current UpCloud account."
}

// return a multiline string man page for the Operation
func (securityUser *UpcloudSecurityUserOperation) Help() string {
	return ""
}

// Is this operation meant to be used only inside the API
func (securityUser *UpcloudSecurityUserOperation) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}

// Run a validation check on the Operation
func (securityUser *UpcloudSecurityUserOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// What settings/values does the Operation provide to an implemenentor
func (securityUser *UpcloudSecurityUserOperation) Properties() api_property.Properties {
	return api_property.New_SimplePropertiesEmpty().Properties()
}

// Execute the Operation
func (securityUser *UpcloudSecurityUserOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	service := securityUser.ServiceWrapper()

	account, err := service.GetAccount()
	if err == nil {
		log.WithFields(log.Fields{"username": account.UserName, "credits": account.Credits}).Info("Current UpCloud Account")
		res.MarkSuccess()
	} else {
		res.AddError(err)
		res.AddError(errors.New("Could not retrieve UpCloud account information."))
		res.MarkFailed()
	}

	res.MarkFinished()

	return res.Result()
}
