package upcloud

import (
	log "github.com/Sirupsen/logrus"

	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"
	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"

	api_property "github.com/wunderkraut/radi-api/property"
	api_usage "github.com/wunderkraut/radi-api/usage"
)

/**
 * Custom properties for the upcloud operations
 */

const (
	UPCLOUD_GLOBAL_PROPERTY               = "upcloud.global"
	UPCLOUD_FORCE_PROPERTY                = "upcloud.force"
	UPCLOUD_WAIT_PROPERTY                 = "upcloud.wait"
	UPCLOUD_FIREWALL_RULES_PROPERTY       = "upcloud.firewall.rules"
	UPCLOUD_SERVER_UUID_PROPERTY          = "upcloud.server.uuid"
	UPCLOUD_SERVER_UUIDS_PROPERTY         = "upcloud.server.uuids"
	UPCLOUD_SERVER_DETAILS_PROPERTY       = "upcloud.server.details"
	UPCLOUD_SERVER_CREATEREQUEST_PROPERTY = "upcloud.server.createrequest"
	UPCLOUD_STORAGE_UUID_PROPERTY         = "upcloud.storage.uuid"
	UPCLOUD_STORAGE_UUIDS_PROPERTY        = "upcloud.storage.uuids"
	UPCLOUD_ZONE_ID_PROPERTY              = "upcloud.zone.id"
)

// A boolean flag that tells upcloud to consider services/zones outside the scope of the project
// @NOTE this is kind of risky to use, so it should be limited to safe operations
type UpcloudGlobalProperty struct {
	api_property.BooleanProperty
}

// ID returns string unique property Identifier
func (global *UpcloudGlobalProperty) Id() string {
	return UPCLOUD_GLOBAL_PROPERTY
}

// Label returns a short user readable label for the property
func (global *UpcloudGlobalProperty) Label() string {
	return "Global UpCloud services"
}

// Description provides a longer multi-line string description of what the property does
func (global *UpcloudGlobalProperty) Description() string {
	return "Consider UpCloud Service/Zones outside the scope of the project"
}

// Mark a property as being for internal use only (no shown to users)
func (global *UpcloudGlobalProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A boolean flag that tells upcloud to force operations to proceed even if
// blocked by server status.  This may require an additional operation process.
// For example, when deleting a running server, this would first stop the server
type UpcloudForceProperty struct {
	api_property.BooleanProperty
}

// ID returns string unique property Identifier
func (force *UpcloudForceProperty) Id() string {
	return UPCLOUD_FORCE_PROPERTY
}

// Label returns a short user readable label for the property
func (force *UpcloudForceProperty) Label() string {
	return "Force"
}

// Description provides a longer multi-line string description of what the property does
func (force *UpcloudForceProperty) Description() string {
	return "Force the operation"
}

// Mark a property as being for internal use only (no shown to users)
func (force *UpcloudForceProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A boolean flag that tells that command to stay attached until the operation is complete
type UpcloudWaitProperty struct {
	api_property.BooleanProperty
}

// ID returns string unique property Identifier
func (wait *UpcloudWaitProperty) Id() string {
	return UPCLOUD_WAIT_PROPERTY
}

// Label returns a short user readable label for the property
func (wait *UpcloudWaitProperty) Label() string {
	return "Wait for UpCloud finish"
}

// Description provides a longer multi-line string description of what the property does
func (wait *UpcloudWaitProperty) Description() string {
	return "Wait for UpCloud to report the desired change of state before disconnecting"
}

// Mark a property as being for internal use only (no shown to users)
func (wait *UpcloudWaitProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A string slice property to match to server UUID
type UpcloudServerUUIDProperty struct {
	api_property.StringProperty
}

// ID returns string unique property Identifier
func (uuid *UpcloudServerUUIDProperty) Id() string {
	return UPCLOUD_SERVER_UUID_PROPERTY
}

// Label returns a short user readable label for the property
func (uuid *UpcloudServerUUIDProperty) Label() string {
	return "UpCloud server UUID"
}

// Description provides a longer multi-line string description of what the property does
func (uuid *UpcloudServerUUIDProperty) Description() string {
	return "Specific UpCloud server UUID"
}

// Mark a property as being for internal use only (no shown to users)
func (uuid *UpcloudServerUUIDProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A string slice property to match to server UUID
type UpcloudServerUUIDSProperty struct {
	api_property.StringSliceProperty
}

// ID returns string unique property Identifier
func (uuids *UpcloudServerUUIDSProperty) Id() string {
	return UPCLOUD_SERVER_UUIDS_PROPERTY
}

// Label returns a short user readable label for the property
func (uuids *UpcloudServerUUIDSProperty) Label() string {
	return "UpCloud server UUID slice"
}

// Description provides a longer multi-line string description of what the property does
func (uuids *UpcloudServerUUIDSProperty) Description() string {
	return "List of UpCloud server UUIDs"
}

// Mark a property as being for internal use only (no shown to users)
func (uuids *UpcloudServerUUIDSProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A string slice property to match to storage UUID
type UpcloudStorageUUIDProperty struct {
	api_property.StringProperty
}

// ID returns string unique property Identifier
func (uuid *UpcloudStorageUUIDProperty) Id() string {
	return UPCLOUD_STORAGE_UUID_PROPERTY
}

// Label returns a short user readable label for the property
func (uuid *UpcloudStorageUUIDProperty) Label() string {
	return "UpCloud storage UUID"
}

// Description provides a longer multi-line string description of what the property does
func (uuid *UpcloudStorageUUIDProperty) Description() string {
	return "Single UpCloud storage UUID"
}

// Mark a property as being for internal use only (no shown to users)
func (uuid *UpcloudStorageUUIDProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A string slice property to match to storage UUID
type UpcloudStorageUUIDSProperty struct {
	api_property.StringSliceProperty
}

// ID returns string unique property Identifier
func (uuids *UpcloudStorageUUIDSProperty) Id() string {
	return UPCLOUD_STORAGE_UUID_PROPERTY
}

// Label returns a short user readable label for the property
func (uuids *UpcloudStorageUUIDSProperty) Label() string {
	return "UpCloud storage UUID slice"
}

// Description provides a longer multi-line string description of what the property does
func (uuids *UpcloudStorageUUIDSProperty) Description() string {
	return "List of UpCloud server UUIDs"
}

// Mark a property as being for internal use only (no shown to users)
func (uuids *UpcloudStorageUUIDSProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A string slice property to match to zone id
type UpcloudZoneIdProperty struct {
	api_property.StringSliceProperty
}

// ID returns string unique property Identifier
func (id *UpcloudZoneIdProperty) Id() string {
	return UPCLOUD_ZONE_ID_PROPERTY
}

// Label returns a short user readable label for the property
func (id *UpcloudZoneIdProperty) Label() string {
	return "UpCloud zone ID"
}

// Description provides a longer multi-line string description of what the property does
func (id *UpcloudZoneIdProperty) Description() string {
	return "Specific UpCloud zone ID"
}

// Mark a property as being for internal use only (no shown to users)
func (id *UpcloudZoneIdProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// A property for the ServerDetails, not really meant for public consumption
type UpcloudServerDetailsProperty struct {
	value upcloud.ServerDetails
}

// ID returns string unique property Identifier
func (details *UpcloudServerDetailsProperty) Id() string {
	return UPCLOUD_SERVER_DETAILS_PROPERTY
}

// Label returns a short user readable label for the property
func (details *UpcloudServerDetailsProperty) Label() string {
	return "UpCloud details"
}

// Description provides a longer multi-line string description of what the property does
func (details *UpcloudServerDetailsProperty) Description() string {
	return "UpCloud server details object"
}

// Mark a property as being for internal use only (no shown to users)
func (details *UpcloudServerDetailsProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// Give an idea of what type of value the property consumes
func (details *UpcloudServerDetailsProperty) Type() string {
	return "github.com/Jalle19/upcloud-go-sdk/upcloud/ServerDetails"
}

func (details *UpcloudServerDetailsProperty) Get() interface{} {
	return interface{}(details.value)
}
func (details *UpcloudServerDetailsProperty) Set(value interface{}) bool {
	if converted, ok := value.(upcloud.ServerDetails); ok {
		details.value = converted
		return true
	} else {
		log.WithFields(log.Fields{"value": value}).Error("Could not assign Property value, because the passed parameter was the wrong type. Expected UpCloud ServerDetails")
		return false
	}
}

// A property for the CreateServerRequest, not really meant for public consumption
type UpcloudServerCreateRequestProperty struct {
	value upcloud_request.CreateServerRequest
}

// ID returns string unique property Identifier
func (request *UpcloudServerCreateRequestProperty) Id() string {
	return UPCLOUD_SERVER_CREATEREQUEST_PROPERTY
}

// Label returns a short user readable label for the property
func (request *UpcloudServerCreateRequestProperty) Label() string {
	return "UpCloud create request"
}

// Description provides a longer multi-line string description of what the property does
func (request *UpcloudServerCreateRequestProperty) Description() string {
	return "UpCloud server create request object"
}

// Mark a property as being for internal use only (no shown to users)
func (request *UpcloudServerCreateRequestProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// Give an idea of what type of value the property consumes
func (request *UpcloudServerCreateRequestProperty) Type() string {
	return "github.com/Jalle19/upcloud-go-sdk/upcloud/request/CreateServerRequest"
}

func (request *UpcloudServerCreateRequestProperty) Get() interface{} {
	return interface{}(request.value)
}
func (request *UpcloudServerCreateRequestProperty) Set(value interface{}) bool {
	if converted, ok := value.(upcloud_request.CreateServerRequest); ok {
		request.value = converted
		return true
	} else {
		log.WithFields(log.Fields{"value": value}).Error("Could not assign Property value, because the passed parameter was the wrong type. Expected UpCloud Request CreateServerRequest")
		return false
	}
}

// A property for the CreateServerRequest, not really meant for public consumption
type UpcloudFirewallRulesProperty struct {
	value upcloud.FirewallRules
}

// ID returns string unique property Identifier
func (firewallRules *UpcloudFirewallRulesProperty) Id() string {
	return UPCLOUD_FIREWALL_RULES_PROPERTY
}

// Label returns a short user readable label for the property
func (firewallRules *UpcloudFirewallRulesProperty) Label() string {
	return "UpCloud firewall rules"
}

// Description provides a longer multi-line string description of what the property does
func (firewallRules *UpcloudFirewallRulesProperty) Description() string {
	return "UpCloud server firewall rules object"
}

// Mark a property as being for internal use only (no shown to users)
func (firewallRules *UpcloudFirewallRulesProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// Give an idea of what type of value the property consumes
func (firewallRules *UpcloudFirewallRulesProperty) Type() string {
	return "github.com/Jalle19/upcloud-go-sdk/upcloud/FirewallRules"
}

func (firewallRules *UpcloudFirewallRulesProperty) Get() interface{} {
	return interface{}(firewallRules.value)
}
func (firewallRules *UpcloudFirewallRulesProperty) Set(value interface{}) bool {
	if converted, ok := value.(upcloud.FirewallRules); ok {
		firewallRules.value = converted
		return true
	} else {
		log.WithFields(log.Fields{"value": value}).Error("Could not assign Property value, because the passed parameter was the wrong type. Expected a UpCloud FirewallRules objects")
		return false
	}
}
