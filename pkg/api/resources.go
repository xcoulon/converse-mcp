package api

func NewResource(name string, uri string) Resource {
	return Resource{
		Name: name,
		Uri:  uri,
	}
}

func (r Resource) WithMetadata(metadata map[string]any) Resource {
	r.Meta = metadata
	return r
}

func (r Resource) WithDescription(description string) Resource {
	r.Description = &description
	return r
}

func (r Resource) WithTitle(title string) Resource {
	r.Title = &title
	return r
}

func (r Resource) WithMimeType(mimeType string) Resource {
	r.MimeType = &mimeType
	return r
}

func (r Resource) WithSize(size int) Resource {
	r.Size = &size
	return r
}
