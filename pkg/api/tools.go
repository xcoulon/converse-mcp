package api

func NewTool(name string) Tool {
	return Tool{
		Name:        name,
		Annotations: &ToolAnnotations{},
		InputSchema: ToolInputSchema{
			Type:       "object",
			Properties: map[string]map[string]any{},
			Required:   []string{},
		},
		OutputSchema: &ToolOutputSchema{
			Type:       "object",
			Properties: map[string]map[string]any{},
		},
	}
}

const String = "string"

func (t Tool) WithInputProperty(name string, propType string, description string, required bool) Tool {
	t.InputSchema.Properties[name] = map[string]any{
		"type":        propType,
		"description": description,
	}
	if required {
		t.InputSchema.Required = append(t.InputSchema.Required, name)
	}
	return t
}

func (t Tool) WithOutputProperty(name string, propType string, description string, required bool) Tool {
	t.OutputSchema.Properties[name] = map[string]any{
		"type":        propType,
		"description": description,
	}
	if required {
		t.OutputSchema.Required = append(t.OutputSchema.Required, name)
	}
	return t
}

func (t Tool) WithMetadata(metadata map[string]any) Tool {
	t.Meta = metadata
	return t
}

func (t Tool) WithAnnotations(annotations ToolAnnotations) Tool {
	t.Annotations = &annotations
	return t
}

func (t Tool) WithDescription(description string) Tool {
	t.Description = &description
	return t
}

func (t Tool) WithTitle(title string) Tool {
	t.Title = &title
	return t
}

func (t Tool) WithDestructiveHint(destructiveHint bool) Tool {
	t.Annotations.DestructiveHint = &destructiveHint
	return t
}

func (t Tool) WithReadOnlyHint(readOnlyHint bool) Tool {
	t.Annotations.ReadOnlyHint = &readOnlyHint
	return t
}

func (t Tool) WithOutputSchema(outputProps ...PropertyDefinition) Tool {
	t.OutputSchema = newToolOutputSchema(outputProps)
	return t
}

type PropertyDefinition struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // TODO: enum
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

func newToolOutputSchema(properties []PropertyDefinition) *ToolOutputSchema {
	s := &ToolOutputSchema{
		Type:       "object",
		Properties: map[string]map[string]any{},
		Required:   []string{},
	}
	for _, p := range properties {
		s.Properties[p.Name] = map[string]any{
			"type":        p.Type,
			"description": p.Description,
		}
		if p.Required {
			s.Required = append(s.Required, p.Name)
		}
	}
	return s
}
