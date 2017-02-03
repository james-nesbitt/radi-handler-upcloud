package upcloud

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	api_operation "github.com/wunderkraut/radi-api/operation"
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
func (security *UpcloudSecurityHandler) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	baseOperation := security.BaseUpcloudServiceOperation()

	ops := api_operation.Operations{}
	ops.Add(api_operation.Operation(&UpcloudSecurityUserOperation{BaseUpcloudServiceOperation: *baseOperation}))
	security.operations = &ops

	return api_operation.Result(result)
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

// Is this operation meant to be used only inside the API
func (securityUser *UpcloudSecurityUserOperation) Internal() bool {
	return false
}

// Run a validation check on the Operation
func (securityUser *UpcloudSecurityUserOperation) Validate() bool {
	return true
}

// What settings/values does the Operation provide to an implemenentor
func (securityUser *UpcloudSecurityUserOperation) Properties() api_operation.Properties {
	return api_operation.Properties{}
}

// Execute the Operation
func (securityUser *UpcloudSecurityUserOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	service := securityUser.ServiceWrapper()

	account, err := service.GetAccount()
	if err == nil {
		log.WithFields(log.Fields{"username": account.UserName, "credits": account.Credits}).Info("Current UpCloud Account")
		result.MarkSuccess()
	} else {
		result.AddError(err)
		result.AddError(errors.New("Could not retrieve UpCloud account information."))
		result.MarkFailed()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
