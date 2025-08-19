package api

var DefaultCapabilities = ServerCapabilities{
	Prompts: &ServerCapabilitiesPrompts{
		ListChanged: ToBoolPtr(false),
	},
	Resources: &ServerCapabilitiesResources{
		ListChanged: ToBoolPtr(false),
	},
	Tools: &ServerCapabilitiesTools{
		ListChanged: ToBoolPtr(false),
	},
}
