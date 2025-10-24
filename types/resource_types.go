package types

// ResourceReader defines a function type for reading resource content.
type ResourceReader func(uri string) ([]OperationContent, error)

// ListResourcesResult represents the result of listing resources.
type ListResourcesResult struct {
	Resources []Resource `json:"resources"`
}

// Resource represents a resource with its metadata and read function.
type Resource struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	URI         string         `json:"uri"` // e.g.: file:///path/to/file.txt
	Read        ResourceReader `json:"-"`
}

// ReadResourceParams represents the parameters for reading a resource.
type ReadResourceParams struct {
	URI string `json:"uri"`
}

// ReadResourceResult represents the result of reading a resource.
type ReadResourceResult struct {
	Content []OperationContent `json:"content"`
}

// NewResource creates a new Resource with the given parameters.
func NewResource(name, description, uri string, reader ResourceReader) *Resource {
	return &Resource{Name: name, Description: description, URI: uri, Read: reader}
}

// NewListResourcesResult creates a new ListResourcesResult with the given resources.
func NewListResourcesResult(resources []Resource) *ListResourcesResult {
	return &ListResourcesResult{Resources: resources}
}

// NewReadResourceParams creates a new ReadResourceParams with the given URI.
func NewReadResourceParams(uri string) *ReadResourceParams {
	return &ReadResourceParams{URI: uri}
}

// NewReadResourceResult creates a new ReadResourceResult with the given content.
func NewReadResourceResult(content []OperationContent) *ReadResourceResult {
	return &ReadResourceResult{Content: content}
}
