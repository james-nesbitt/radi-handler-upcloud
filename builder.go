package upcloud

import (
	log "github.com/Sirupsen/logrus"

	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_result "github.com/wunderkraut/radi-api/result"
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

// Return a string identifier for the Handler (not functionally needed yet)
func (builder *UpcloudBuilder) Id() string {
	return "upcloud"
}

// Set a API for this Handler
func (builder *UpcloudBuilder) SetAPI(parent api_api.API) {
	// Keep that api, so that we can use it to make a ConfigWrapper later on
	builder.parent = parent
}

// Initialize and activate the Handler
func (builder *UpcloudBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) api_result.Result {
	if &builder.handlers == nil {
		builder.handlers = api_handler.Handlers{}
	}

	// process and merge the settings
	settings := UpcloudBuilderSettings{}
	settingsProvider.AssignSettings(&settings)
	builder.settings.Merge(settings)

	// This base handler is commonly used in the implementation handlers, so get it once here.
	baseHandler := builder.base_BaseUpcloudServiceHandler()
	// Start a collective result
	res := api_result.New_StandardResult()

	for _, implementation := range implementations.Order() {
		var handler api_handler.Handler

		switch implementation {
		case "monitor":
			handler = api_handler.Handler(&UpcloudMonitorHandler{BaseUpcloudServiceHandler: *baseHandler})
		case "server":
			handler = api_handler.Handler(&UpcloudServerHandler{BaseUpcloudServiceHandler: *baseHandler})
		case "provision":
			handler = api_handler.Handler(&UpcloudProvisionHandler{BaseUpcloudServiceHandler: *baseHandler})
		case "security":
			handler = api_handler.Handler(&UpcloudSecurityHandler{BaseUpcloudServiceHandler: *baseHandler})
		default:
			log.WithFields(log.Fields{"implementation": implementation}).Error("Unknown implementation in UpCloud builder")
		}

		if handler != nil {
			initRes := handler.Validate()
			<-initRes.Finished()

			if initRes.Success() {
				builder.handlers.Add(api_handler.Handler(handler))
			}

			res.Merge(initRes)
		}

	}

	return nil
}

// Validate the builder after Activation is complete
func (builder *UpcloudBuilder) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Return a list of Operations from the Handler
func (builder *UpcloudBuilder) Operations() api_operation.Operations {
	return builder.handlers.Operations()
}

// Return a shared BaseUpcloudServiceOperation for any operation that needs it
func (builder *UpcloudBuilder) base_BaseUpcloudServiceHandler() *BaseUpcloudServiceHandler {
	if builder.base_UpcloudServiceHandler == nil {
		// Builder a configwrapper, which will be used to build upcloud service structs
		ops := builder.parent.Operations()
		configWrapper := api_config.New_SimpleConfigWrapper(ops)
		// get an upcloud factory, using the config wrapper (probably a file like upcloud.yml)
		upcloudFactory := New_UpcloudFactoryConfigWrapperYaml(configWrapper)
		upcloudFactory.Load()

		// Builder the base operation, and keep it
		builder.base_UpcloudServiceHandler = New_BaseUpcloudServiceHandler(upcloudFactory.UpcloudFactory(), &builder.settings)
	}
	return builder.base_UpcloudServiceHandler
}
