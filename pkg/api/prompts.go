package api

func NewPrompt(name string) Prompt {
	return Prompt{
		Name: name,
	}
}

func (p Prompt) WithMetadata(metadata map[string]any) Prompt {
	p.Meta = metadata
	return p
}

func (p Prompt) WithDescription(description string) Prompt {
	p.Description = &description
	return p
}

func (p Prompt) WithTitle(title string) Prompt {
	p.Title = &title
	return p
}

func (p Prompt) WithArguments(arguments ...PromptArgument) Prompt {
	p.Arguments = arguments
	return p
}
