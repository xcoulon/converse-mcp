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

type ServerCapability func(*ServerCapabilities)

func PromptListChanged(v bool) ServerCapability {
	return func(sc *ServerCapabilities) {
		sc.Prompts.ListChanged = ToBoolPtr(v)
	}
}

func ResourceListChanged(v bool) ServerCapability {
	return func(sc *ServerCapabilities) {
		sc.Resources.ListChanged = ToBoolPtr(v)
	}
}

func ToolListChanged(v bool) ServerCapability {
	return func(sc *ServerCapabilities) {
		sc.Tools.ListChanged = ToBoolPtr(v)
	}
}
