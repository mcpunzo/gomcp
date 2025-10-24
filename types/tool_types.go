package types

// ToolHandler defines the function signature for tool execution handlers.
// It takes an input a struct containing the parameters of the function.
type ToolHandler func(map[string]any) (*ToolResult, error)

// Tool represents a tool that can be called via the MCP protocol.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
	Run         ToolHandler    `json:"-"`
}

// ListToolsResult represents the result of listing available tools.
type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

// CallToolParams represents the parameters for calling a tool.
type CallToolParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// ToolResult represents the result returned by a tool execution.
type ToolResult struct {
	Content []OperationContent `json:"content"`
}

// NewTool creates a new Tool with the given parameters.
func NewTool(name, description string, inputSchema map[string]any, handler ToolHandler) *Tool {
	return &Tool{name, description, inputSchema, handler}
}

// NewListToolsResult creates a new ListToolsResult with the given tools.
func NewListToolsResult(tools []Tool) *ListToolsResult {
	return &ListToolsResult{tools}
}

// NewCallToolParams creates a new CallToolParams with the given name and arguments.
func NewCallToolParams(name string, arguments map[string]any) *CallToolParams {
	return &CallToolParams{name, arguments}
}

// NewToolResult creates a new ToolResult with the given content.
func NewToolResult(content []OperationContent) *ToolResult {
	return &ToolResult{content}
}

// NewOperationContent creates a new OperationContent with the given parameters.
func NewOperationContent(_type, text, uri string, data any) *OperationContent {
	return &OperationContent{_type, text, data, uri}
}
