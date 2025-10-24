package gomcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/mcpunzo/gomcp/internal/type_converter"
	"github.com/mcpunzo/gomcp/types"
)

const (
	ErrParse          = -32700
	ErrInvalidRequest = -32600
	ErrMethodNotFound = -32601
	ErrInvalidParams  = -32602
	ErrInternal       = -32603

	ErrServerGeneric = -32000
	ErrAccessDenied  = -32001
	ErrNotFound      = -32002
)

const (
	Initialize    = "initialize"
	Shutdown      = "shutdown"
	ListTools     = "tools/list"
	CallTool      = "tools/call"
	ListResources = "resources/list"
	ReadResource  = "resources/read"
)

type MCPServer struct {
	name            string
	version         string
	shutdownMessage string
	tools           map[string]*types.Tool
	resources       map[string]*types.Resource
}

// New creates a new MCPServer instance with the given name and version.
func New(name, version string) *MCPServer {
	return &MCPServer{name, version, "MCP Session terminated", make(map[string]*types.Tool), make(map[string]*types.Resource)}
}

// AddTool adds a tool to the MCPServer.
func (m *MCPServer) AddTool(tool *types.Tool) {
	m.tools[tool.Name] = tool
}

func (m *MCPServer) AddToolFunc(name, description string, handler any) error {
	handlerType := reflect.TypeOf(handler)

	// handler must be a func
	if handlerType.Kind() != reflect.Func {
		return errors.New("handler must be a function")
	}

	// check the func signature
	if handlerType.NumIn() != 1 {
		return errors.New("handler must accept exactly 1 argument")
	}

	if handlerType.NumOut() != 2 {
		return errors.New("handler must return exactly 2 values")
	}

	// input type must be a struct
	argType := handlerType.In(0)
	if argType.Kind() != reflect.Struct {
		return errors.New("handler argument must be a struct")
	}

	// check the output types
	if handlerType.Out(0) != reflect.TypeOf((*types.ToolResult)(nil)) {
		return errors.New("handler first return value must be *types.ToolResult")
	}
	if !handlerType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.New("handler second return value must be error")
	}

	// Generate the InputSchema from the struct
	schema := m.generateJSONSchema(argType)

	// Create a wrapper to convert map[string]interface{} -> struct
	wrappedHandler := func(args map[string]interface{}) (*types.ToolResult, error) {
		// Serializza args in JSON
		jsonData, err := json.Marshal(args)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal args: %w", err)
		}

		// Create a new  input struct instance
		argValue := reflect.New(argType).Interface()

		// Deserialize JSON into the struct
		if err := json.Unmarshal(jsonData, argValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal args: %w", err)
		}

		// Invoke the original handler
		handlerValue := reflect.ValueOf(handler)
		results := handlerValue.Call([]reflect.Value{
			reflect.ValueOf(argValue).Elem(),
		})

		// Extract results
		var result *types.ToolResult
		var errResult error

		if !results[0].IsNil() {
			result = results[0].Interface().(*types.ToolResult)
		}
		if !results[1].IsNil() {
			errResult = results[1].Interface().(error)
		}

		return result, errResult
	}

	m.AddTool(types.NewTool(name, description, schema, wrappedHandler))

	return nil
}

func (m *MCPServer) generateJSONSchema(t reflect.Type) map[string]any {
	props := map[string]any{}
	required := []string{}

	for i := range t.NumField() {
		field := t.Field(i)
		name := field.Name
		props[name] = map[string]any{"type": field.Type.Kind().String()}
		required = append(required, name)
	}

	inputSchema := map[string]any{
		"type":       "object",
		"properties": props,
		"required":   required,
	}

	return inputSchema
}

// AddResource adds a resource to the MCPServer.
func (m *MCPServer) AddResource(resource *types.Resource) {
	m.resources[resource.URI] = resource
}

// Tools returns a list of all registered tools.
func (m *MCPServer) Tools() []types.Tool {
	return type_converter.MapValueToArray(m.tools)
}

// Resources returns a list of all registered resources.
func (m *MCPServer) Resources() []types.Resource {
	return type_converter.MapValueToArray(m.resources)
}

// HandleRequest handles an incoming JSON-RPC request and returns the appropriate response.
func (m *MCPServer) HandleRequest(req *types.JSONRPCRequest) *types.JSONRPCResponse {
	log.Printf("Handling request: %s", req.Method)

	switch req.Method {
	case Initialize:
		return m.handleInitialize(req)
	case Shutdown:
		return m.handleShutdown(req)
	case ListTools:
		return m.handleListTools(req)
	case CallTool:
		return m.handleCallTool(req)
	case ListResources:
		return m.handleListResources(req)
	case ReadResource:
		return m.handleReadResource(req)
	default:
		return m.handleError(req.Id, "Method Not Found", ErrMethodNotFound, req.Method)
	}
}

func (m *MCPServer) handleInitialize(req *types.JSONRPCRequest) *types.JSONRPCResponse {
	return types.NewJSONRPCResponse(req.Id, types.NewInitializeResult(m.name, m.version, len(m.tools) > 0, len(m.resources) > 0), nil)
}

func (m *MCPServer) handleError(id, message string, code int, data any) *types.JSONRPCResponse {
	return types.NewJSONRPCResponse(id, nil, types.NewJSONRPCErrorObj(code, message, data))
}

func (m *MCPServer) handleShutdown(req *types.JSONRPCRequest) *types.JSONRPCResponse {
	return types.NewJSONRPCResponse(req.Id, types.NewShutdownResult(m.shutdownMessage), nil)
}

func (m *MCPServer) handleListTools(req *types.JSONRPCRequest) *types.JSONRPCResponse {
	return types.NewJSONRPCResponse(req.Id, types.NewListToolsResult(m.Tools()), nil)
}

func (m *MCPServer) handleListResources(req *types.JSONRPCRequest) *types.JSONRPCResponse {
	return types.NewJSONRPCResponse(req.Id, types.NewListResourcesResult(m.Resources()), nil)
}

func (m *MCPServer) handleCallTool(req *types.JSONRPCRequest) *types.JSONRPCResponse {
	params, ok := req.Params.(*types.CallToolParams)

	if !ok {
		return m.handleError(req.Id, "Invalid parameters", ErrInvalidParams, req.Method)
	}

	tool, exists := m.tools[params.Name]
	if !exists {
		return m.handleError(req.Id, "Unknown Tool", ErrMethodNotFound, req.Method)
	}

	res, err := tool.Run(params.Arguments)
	if err != nil {
		return m.handleError(req.Id, fmt.Sprintf("Error executing tool %v", tool.Name), ErrServerGeneric, err.Error())
	}

	return types.NewJSONRPCResponse(req.Id, res, nil)
}

func (m *MCPServer) handleReadResource(req *types.JSONRPCRequest) *types.JSONRPCResponse {
	params, ok := req.Params.(*types.ReadResourceParams)

	if !ok {
		return m.handleError(req.Id, "Invalid parameters", ErrInvalidParams, req.Method)
	}

	resource, exists := m.resources[params.URI]
	if !exists {
		return m.handleError(req.Id, "Unknown Resource", ErrMethodNotFound, req.Method)
	}

	content, err := resource.Read(params.URI)
	if err != nil {
		return m.handleError(req.Id, fmt.Sprintf("Error reading resource %v", resource.Name), ErrServerGeneric, err.Error())
	}

	return types.NewJSONRPCResponse(req.Id, types.NewReadResourceResult(content), nil)
}
