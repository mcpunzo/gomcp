package types

// ClientInfo represents information about the client.
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerInfo represents information about the server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Capabilities represents the server's capabilities.
type Capabilities struct {
	Tools     bool `json:"tools"`
	Resources bool `json:"resources"`
}

// InitializeParams represents the parameters for the initialize request.
type InitializeParams struct {
	ClientInfo ClientInfo `json:"clientInfo"`
}

// InitializeResult represents the result of the initialize request.
type InitializeResult struct {
	ServerInfo   ServerInfo   `json:"serverInfo"`
	Capabilities Capabilities `json:"capabilities"`
}

// NewInitializeParams creates a new InitializeParams instance.
func NewInitializeParams(name, version string) *InitializeParams {
	return &InitializeParams{ClientInfo: ClientInfo{Name: name, Version: version}}
}

// NewInitializeResult creates a new InitializeResult instance.
func NewInitializeResult(name, version string, tools, resources bool) *InitializeResult {
	return &InitializeResult{
		ServerInfo:   ServerInfo{Name: name, Version: version},
		Capabilities: Capabilities{Tools: tools, Resources: resources},
	}
}
