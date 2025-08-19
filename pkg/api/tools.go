package api

func NewTool(name string, inputProps ...PropertyDefinition) Tool {
	return Tool{
		Name:        name,
		InputSchema: newToolInputSchema(inputProps),
	}
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

func newToolInputSchema(properties []PropertyDefinition) ToolInputSchema {
	s := ToolInputSchema{
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
