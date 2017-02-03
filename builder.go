package upcloud

import (
	log "github.com/Sirupsen/logrus"

	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"
)

/**
 * A radi builder for upcloud handlers
 */

// Upcloud Builder
type UpcloudBuilder struct {
	parent   api_api.API
	handlers api_handler.Handlers

	settings UpcloudBuilderSettings

	base_UpcloudServiceHandler *BaseUpcloudServiceHandler
}

// Set a API for this Handler
func (builder *UpcloudBuilder) SetAPI(parent api_api.API) {
	// Keep that api, so that we can use it to make a ConfigWrapper later on
	builder.parent = parent
}

// Initialize and activate the Handler
func (builder *UpcloudBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) error {
	if &builder.handlers == nil {
		builder.handlers = api_handler.Handlers{}
	}

	// process and merge the settings
	settings := UpcloudBuilderSettings{}
	settingsProvider.AssignSettings(&settings)
	builder.settings.Merge(settings)

	// This base handler is commonly used in the implementation handlers, so get it once here.
	baseHandler := builder.base_BaseUpcloudServiceHandler()

	for _, implementation := range implementations.Order() {
		switch implementation {
		case "monitor":
			monitorHandler := UpcloudMonitorHandler{BaseUpcloudServiceHandler: *baseHandler}
			monitorHandler.Init()
			builder.handlers.Add(api_handler.Handler(&monitorHandler))
		case "server":
			serverHandler := UpcloudServerHandler{BaseUpcloudServiceHandler: *baseHandler}
			serverHandler.Init()
			builder.handlers.Add(api_handler.Handler(&serverHandler))
		case "provision":
			provisionHandler := UpcloudProvisionHandler{BaseUpcloudServiceHandler: *baseHandler}
			provisionHandler.Init()
			builder.handlers.Add(api_handler.Handler(&provisionHandler))
		case "security":
			securityHandler := UpcloudSecurityHandler{BaseUpcloudServiceHandler: *baseHandler}
			securityHandler.Init()
			builder.handlers.Add(api_handler.Handler(&securityHandler))
		default:
			log.WithFields(log.Fields{"implementation": implementation}).Error("Unknown implementation in UpCloud builder")
		}
	}

	return nil
}

// Rturn a string identifier for the Handler (not functionally needed yet)
func (builder *UpcloudBuilder) Id() string {
	return "upcloud"
}

// Return a list of Operations from the Handler
func (builder *UpcloudBuilder) Operations() *api_operation.Operations {
	ops := builder.handlers.Operations()
	return &ops
}

// Return a shared BaseUpcloudServiceOperation for any operation that needs it
func (builder *UpcloudBuilder) base_BaseUpcloudServiceHandler() *BaseUpcloudServiceHandler {
	if builder.base_UpcloudServiceHandler == nil {
		// Builder a configwrapper, which will be used to build upcloud service structs
		operations := builder.parent.Operations()
		configWrapper := api_config.New_SimpleConfigWrapper(&operations)
		// get an upcloud factory, using the config wrapper (probably a file like upcloud.yml)
		upcloudFactory := New_UpcloudFactoryConfigWrapperYaml(configWrapper)
		upcloudFactory.Load()

		// Builder the base operation, and keep it
		builder.base_UpcloudServiceHandler = New_BaseUpcloudServiceHandler(upcloudFactory.UpcloudFactory(), &builder.settings)
	}
	return builder.base_UpcloudServiceHandler
}
