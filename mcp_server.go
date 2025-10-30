package gomcp

import (
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

// AddToolFunc adds a tool to the MCPServer using a simplified handler function.
// The handler function should have the signature:
//
//	func(args AnyStruct) (*types.ToolResult, error)
//
// where AnyStruct is a struct type defining the input parameters for the tool.
// The input schema is automatically generated based on the fields of the struct.
func (m *MCPServer) AddToolFunc(name, description string, handler any) {
	v := reflect.ValueOf(handler)
	if v.Kind() != reflect.Func {
		panic("handler must be a function func (args AnyStruct) (*types.ToolResult, error)")
	}

	t := v.Type()
	if t.NumIn() == 0 {
		panic("handler need to have a parameter of type struct")
	}

	arg := t.In(0)
	if arg.Kind() == reflect.Ptr {
		arg = arg.Elem()
	}

	if arg.Kind() != reflect.Struct {
		panic("handler parameter need to be of type struct")
	}

	props := map[string]any{}
	required := []string{}

	for i := range arg.NumField() {
		field := arg.Field(i)
		name := field.Name
		props[name] = map[string]any{"type": field.Type.Kind().String()}
		required = append(required, name)
	}

	inputSchema := map[string]any{
		"type":       "object",
		"properties": props,
		"required":   required,
	}

	m.tools[name] = types.NewTool(name, description, inputSchema, func(args any) (*types.ToolResult, error) {
		v := reflect.ValueOf(handler)
		in := []reflect.Value{reflect.ValueOf(args)}

		results := v.Call(in)

		if len(results) != 2 {
			return nil, fmt.Errorf("handler must return (*types.ToolResult, error)")
		}

		var res *types.ToolResult
		var err error

		// ✅ Safe parse primo valore
		rv := results[0]
		if rv.IsValid() {
			// Non chiamare IsNil se non è nillable!
			if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface ||
				rv.Kind() == reflect.Map || rv.Kind() == reflect.Slice ||
				rv.Kind() == reflect.Func || rv.Kind() == reflect.Chan {
				if rv.IsNil() {
					// ok, lascia res = nil
				} else {
					val := rv.Interface()
					if r, ok := val.(*types.ToolResult); ok {
						res = r
					} else {
						return nil, fmt.Errorf("handler must return (*types.ToolResult, error)")
					}
				}
			} else {
				// tipo non nillable → errore
				return nil, fmt.Errorf("handler must return (*types.ToolResult, error)")
			}
		}

		// ✅ Safe parse secondo valore
		ev := results[1]
		if ev.IsValid() {
			if ev.Kind() == reflect.Interface || ev.Kind() == reflect.Ptr {
				if !ev.IsNil() {
					val := ev.Interface()
					if e, ok := val.(error); ok {
						err = e
					} else {
						return nil, fmt.Errorf("handler must return (*types.ToolResult, error)")
					}
				}
			} else if ev.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				err, _ = ev.Interface().(error)
			} else {
				return nil, fmt.Errorf("handler must return (*types.ToolResult, error)")
			}
		}

		return res, err
	})
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
