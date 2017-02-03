package upcloud

import (
	log "github.com/Sirupsen/logrus"

	upcloud "github.com/Jalle19/upcloud-go-sdk/upcloud"
	upcloud_request "github.com/Jalle19/upcloud-go-sdk/upcloud/request"

	api_operation "github.com/wunderkraut/radi-api/operation"
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
	api_operation.BooleanProperty
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
func (global *UpcloudGlobalProperty) Internal() bool {
	return false
}

// A boolean flag that tells upcloud to force operations to proceed even if
// blocked by server status.  This may require an additional operation process.
// For example, when deleting a running server, this would first stop the server
type UpcloudForceProperty struct {
	api_operation.BooleanProperty
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
func (force *UpcloudForceProperty) Internal() bool {
	return false
}

// A boolean flag that tells that command to stay attached until the operation is complete
type UpcloudWaitProperty struct {
	api_operation.BooleanProperty
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
func (wait *UpcloudWaitProperty) Internal() bool {
	return false
}

// A string slice property to match to server UUID
type UpcloudServerUUIDProperty struct {
	api_operation.StringProperty
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
func (uuid *UpcloudServerUUIDProperty) Internal() bool {
	return false
}

// A string slice property to match to server UUID
type UpcloudServerUUIDSProperty struct {
	api_operation.StringSliceProperty
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
func (uuids *UpcloudServerUUIDSProperty) Internal() bool {
	return false
}

// A string slice property to match to storage UUID
type UpcloudStorageUUIDProperty struct {
	api_operation.StringProperty
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
func (uuid *UpcloudStorageUUIDProperty) Internal() bool {
	return false
}

// A string slice property to match to storage UUID
type UpcloudStorageUUIDSProperty struct {
	api_operation.StringSliceProperty
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
func (uuids *UpcloudStorageUUIDSProperty) Internal() bool {
	return false
}

// A string slice property to match to zone id
type UpcloudZoneIdProperty struct {
	api_operation.StringSliceProperty
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
func (id *UpcloudZoneIdProperty) Internal() bool {
	return false
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
func (details *UpcloudServerDetailsProperty) Internal() bool {
	return false
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
func (request *UpcloudServerCreateRequestProperty) Internal() bool {
	return false
}

// Give an idea of what type of value the property consumes
func (request *UpcloudServerCreateRequestProperty) Type() string {
	return "github.com/Jalle19/upcloud-go-sdk/upcloud/request/CreateServerRequest"
}

func (request *UpcloudServerCreateRequestProperty) Get() interface{} {
	// // prov_project := "raditest"
	// prov_zone := "fi-hel1"

	// prov_initscript := "" // Initialize script. Can be a URL

	// prov_user := upcloud_request.LoginUser{
	// 	CreatePassword: "no", // Allow SSH only with key
	// 	Username:       "radi",
	// 	SSHKeys: []string{
	// 		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDEI7j4EaK2RRKgp7rA9gDIL279WtNBWsPQwKn6YNjb7i1EUAM+IYzdQbPgpYr0rMx67DhvbK1pBeL0HTXfk1ZnSbbZe2xktk+YJo6l8zQ7wYydWMcCcB5HUvgG1/ugTj6wxImYAx7sEuXY4MVO7aHmfMnjV+7Re0uXHjAPL9k5O2Xvt75RmrgG8YpE6MvZtYTzIRmINbuSAX9CWKi46ZuRNYKDyZTSarqA1TOaGKiO6vf2dM7bWqsvitOxwEC6Z0c5nIAjcAGhg+yBEEloWTqNqkxPzbh0AIIO9HjGlnbSaIffwrv78UzrHatukUQOcsg6PBvMPvhjdoB0JrscLneDy0DhF6ptAQporg3SieypB/3hiZ0RfT94c35DQufFphfsphIBXIsqENJKR383sz57PPDtVgyXKSu5ujhXUPgC1rwldGUqVtMLsvC4tcnOIbOK917QjUQ+8cJoC08XBUG10knUoIWP8ebv55sfnBTHW27g+4B1V6ub3Zyq/ZRzeJXWzSh1QmOUXp1Q47bEz3eT2/VRtKyUYAo3ChvceMSERsVtwfRgIcAreIqGy2GJQPe7NkYOAPrirwhfppoFJ/nx3xGFjg0iZIg4Z1nTpwEWBWcC5eo/yfORnhdAooJWRYO37nOrjryUZJsRbFC/Uj7JOIX2QrZEX1bm4SwgSF8tTQ==", // JRN
	// 	},
	// }
	// prov_networks := []upcloud_request.CreateServerIPAddress{
	// 	upcloud_request.CreateServerIPAddress{
	// 		Access: "private",
	// 		Family: "IPv4",
	// 	},
	// 	upcloud_request.CreateServerIPAddress{
	// 		Access: "public",
	// 		Family: "IPv4",
	// 	},
	// 	upcloud_request.CreateServerIPAddress{
	// 		Access: "public",
	// 		Family: "IPv6",
	// 	},
	// }
	// prov_storages := []upcloud.CreateServerStorageDevice{
	// 	upcloud.CreateServerStorageDevice{ // primary disk
	// 		Action:  "clone",
	// 		Storage: "01000000-0000-4000-8000-000080010200", //size=5 title="CoreOS Stable 1068.8.0" type=template
	// 		Title:   "coreos-install",
	// 		Size:    10, // Storage size in gigabytes, if cloning it has to be larger the source size
	// 		Tier:    "maxiops",
	// 	},
	// 	upcloud.CreateServerStorageDevice{ // primary disk
	// 		Action:  "create",
	// 		Address: "virtio:0",
	// 		Title:   "coreos-root",
	// 		Size:    10,
	// 		Tier:    "maxiops",
	// 		Type:    "disk",
	// 	},
	// }
	// prov_firewall_rules := []upcloud.FirewallRule{
	// 	// upcloud.FirewallRule{
	// 	// 	Action:                  "accept",
	// 	// 	Comment:                 "Alow HTTP from anywhere",
	// 	// 	DestinationAddressEnd:   "",
	// 	// 	DestinationAddressStart: "",
	// 	// 	DestinationPortEnd:      "80",
	// 	// 	DestinationPortStart:    "80",
	// 	// 	Direction:               "in",
	// 	// 	Family:                  "IPv4",
	// 	// 	ICMPType:                "",
	// 	// 	Position:                1,
	// 	// 	Protocol:                "",
	// 	// 	SourceAddressEnd:        "",
	// 	// 	SourceAddressStart:      "",
	// 	// 	SourcePortEnd:           "",
	// 	// 	SourcePortStart:         "",
	// 	// },
	// 	upcloud.FirewallRule{
	// 		Action:               "accept",
	// 		Comment:              "Alow HTTP from anywhere",
	// 		DestinationPortEnd:   "80",
	// 		DestinationPortStart: "80",
	// 		Direction:            "in",
	// 		Family:               "IPv4",
	// 		Position:             1,
	// 		Protocol:             "tcp",
	// 	},
	// 	upcloud.FirewallRule{
	// 		Action:               "accept",
	// 		Comment:              "Allow SSH from a specific network only",
	// 		DestinationPortEnd:   "22",
	// 		DestinationPortStart: "22",
	// 		Direction:            "in",
	// 		Family:               "IPv4",
	// 		ICMPType:             "tcp",
	// 		Position:             2,
	// 		Protocol:             "tcp",
	// 		SourceAddressEnd:     "192.168.1.255",
	// 		SourceAddressStart:   "192.168.1.1",
	// 	},
	// 	upcloud.FirewallRule{
	// 		Action:                  "accept",
	// 		Comment:                 "Allow SSH over IPv6 from this range",
	// 		DestinationAddressEnd:   "",
	// 		DestinationAddressStart: "",
	// 		DestinationPortEnd:      "22",
	// 		DestinationPortStart:    "22",
	// 		Direction:               "in",
	// 		Family:                  "IPv6",
	// 		ICMPType:                "",
	// 		Position:                3,
	// 		Protocol:                "tcp",
	// 		SourceAddressEnd:        "2a04:3540:1000:aaaa:bbbb:cccc:d001",
	// 		SourceAddressStart:      "2a04:3540:1000:aaaa:bbbb:cccc:d001",
	// 		SourcePortEnd:           "",
	// 		SourcePortStart:         "",
	// 	},
	// 	upcloud.FirewallRule{
	// 		Action:    "accept",
	// 		Comment:   "Allow ICMP echo request (ping)",
	// 		Direction: "in",
	// 		Family:    "IPv4",
	// 		Position:  4,
	// 		Protocol:  "icmp",
	// 		ICMPType:  "8",
	// 	},
	// 	upcloud.FirewallRule{
	// 		Action:    "drop",
	// 		Direction: "in",
	// 		Position:  5,
	// 	},
	// }

	// // hardcoded_tag := "radi-provisioned-" + prov_project

	// request.value = upcloud_request.CreateServerRequest{
	// 	//AvoidHost  string `xml:"avoid_host,omitempty"`
	// 	//BootOrder  string `xml:"boot_order,omitempty"`
	// 	//CoreNumber int    `xml:"core_number,omitempty"`
	// 	Firewall:    "on",
	// 	Hostname:    prov_project,
	// 	IPAddresses: prov_networks,
	// 	LoginUser:   &prov_user,
	// 	// MemoryAmount: 2048,
	// 	PasswordDelivery: "none",
	// 	Plan:             "1xCPU-1GB",
	// 	StorageDevices:   prov_storages,
	// 	TimeZone:         "Europe/Helsinki",
	// 	Title:            prov_project + ": provisioned automatically by radi",
	// 	UserData:         prov_initscript,
	// 	//VNC: "off",
	// 	Zone: prov_zone,
	// }

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
func (firewallRules *UpcloudFirewallRulesProperty) Internal() bool {
	return false
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
