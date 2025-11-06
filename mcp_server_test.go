package gomcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/mcpunzo/gomcp/types"
)

func setupTest(tb testing.TB) (*MCPServer, func(tb testing.TB)) {
	log.Println("setup test")

	mcpserver := New("serverName", "v1.0")

	if mcpserver == nil {
		tb.Errorf("Expected !nil but was %v", mcpserver)
	}

	return mcpserver, func(tb testing.TB) {
		log.Println("teardown test")
	}
}

func TestMCPServerNew(t *testing.T) {
	_, teardown := setupTest(t)
	defer teardown(t)
}

type MockTransport struct{}

func (m *MockTransport) SetMCPServer(mcpserver *MCPServer) {}
func (m *MockTransport) Start()                            {}

func TestWithTransport(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	mockTransport := &MockTransport{}

	mcpserver.WithTransport(mockTransport)

	if mcpserver.transport != mockTransport {
		t.Errorf("Expected %v but got %v", mockTransport, mcpserver.transport)
	}
}

func TestHandleRequest(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	table := []struct {
		req      *types.JSONRPCRequest
		expected *types.JSONRPCResponse
	}{
		{
			types.NewJSONRPCRequest("id", Initialize, types.NewInitializeParams("test", "1.0")),
			types.NewJSONRPCResponse("id", types.NewInitializeResult(mcpserver.name, mcpserver.version, len(mcpserver.tools) > 0, len(mcpserver.resources) > 0), nil),
		},
		{
			types.NewJSONRPCRequest("id", "UnknownMethod", nil),
			types.NewJSONRPCResponse("id", nil, types.NewJSONRPCErrorObj(ErrMethodNotFound, "Method Not Found", "UnknownMethod")),
		},
		{
			types.NewJSONRPCRequest("id", Shutdown, types.NewShutdownParams()),
			types.NewJSONRPCResponse("id", types.NewShutdownResult(ShutdownMessage), nil),
		},
		{
			types.NewJSONRPCRequest("id", ListTools, nil),
			types.NewJSONRPCResponse("id", types.NewListToolsResult(mcpserver.Tools()), nil),
		},
		{
			types.NewJSONRPCRequest("id", ListResources, nil),
			types.NewJSONRPCResponse("id", types.NewListResourcesResult(mcpserver.Resources()), nil),
		},
	}

	for _, test := range table {
		actual := mcpserver.HandleRequest(test.req)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %#v but got %#v", test.expected, actual)
		}
	}
}

func TestHandleRequestWithCallTool(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	err := errors.New("worng parameter")
	content := []types.OperationContent{*types.NewOperationContent("text", "content of file£", "", nil)}

	tool := types.NewTool(
		"read_file",
		"read_file tool description",
		map[string]any{
			"type":       "object",
			"properties": map[string]any{"path": map[string]string{"type": "string"}},
			"required":   []string{"path"},
		},
		func(args map[string]any) (*types.ToolResult, error) {
			if args == nil {
				return nil, err
			}

			path, ok := args["path"].(string)
			if !ok || path == "" {
				return nil, err
			}

			return types.NewToolResult(content), nil
		},
	)

	mcpserver.AddTool(tool)

	table := []struct {
		req      *types.JSONRPCRequest
		expected *types.JSONRPCResponse
	}{
		{
			types.NewJSONRPCRequest("id", CallTool, types.NewCallToolParams("read_file", map[string]any{"path": "/tmp/example.txt"})),
			types.NewJSONRPCResponse("id", types.NewToolResult(content), nil),
		},
		{
			types.NewJSONRPCRequest("id", CallTool, types.NewCallToolParams("read_file", map[string]any{"wrong_param": "/tmp/example.txt"})),
			types.NewJSONRPCResponse("id", nil, types.NewJSONRPCErrorObj(ErrServerGeneric, "Error executing tool read_file", err.Error())),
		},
		{
			types.NewJSONRPCRequest("id", CallTool, types.NewCallToolParams("not_existing_tool", map[string]any{"path": "/tmp/example.txt"})),
			types.NewJSONRPCResponse("id", nil, types.NewJSONRPCErrorObj(ErrMethodNotFound, "Unknown Tool", CallTool)),
		},

		{
			types.NewJSONRPCRequest("id", CallTool, types.NewShutdownParams()),
			types.NewJSONRPCResponse("id", nil, types.NewJSONRPCErrorObj(ErrInvalidParams, "Invalid parameters", CallTool)),
		},
	}

	for _, test := range table {
		actual := mcpserver.HandleRequest(test.req)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %#v but got %#v", test.expected, actual)
		}
	}
}

func TestHandleRequestWithReadResource(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	content := []types.OperationContent{*types.NewOperationContent("text", "content", "", nil)}

	resource := types.NewResource("res1", "descr1", "file://resource.tst", func(uri string) ([]types.OperationContent, error) {
		return content, nil
	})

	err := errors.New("read error")
	error_resource := types.NewResource("error_resource", "resource generating error", "file://error_resource", func(uri string) ([]types.OperationContent, error) {
		return nil, err
	})

	mcpserver.AddResource(resource)
	mcpserver.AddResource(error_resource)

	table := []struct {
		req      *types.JSONRPCRequest
		expected *types.JSONRPCResponse
	}{
		{
			types.NewJSONRPCRequest("id", ReadResource, types.NewReadResourceParams("file://resource.tst")),
			types.NewJSONRPCResponse("id", types.NewReadResourceResult(content), nil),
		},
		{
			types.NewJSONRPCRequest("id", ReadResource, types.NewReadResourceParams("file://unknown.txt")),
			types.NewJSONRPCResponse("id", nil, types.NewJSONRPCErrorObj(ErrMethodNotFound, "Unknown Resource", ReadResource)),
		},
		{
			types.NewJSONRPCRequest("id", ReadResource, types.NewShutdownParams()),
			types.NewJSONRPCResponse("id", nil, types.NewJSONRPCErrorObj(ErrInvalidParams, "Invalid parameters", ReadResource)),
		},
		{
			types.NewJSONRPCRequest("id", ReadResource, types.NewReadResourceParams("file://error_resource")),
			types.NewJSONRPCResponse("id", nil, types.NewJSONRPCErrorObj(ErrServerGeneric, "Error reading resource error_resource", err.Error())),
		},
	}

	for _, test := range table {
		actual := mcpserver.HandleRequest(test.req)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected %#v but got %#v", test.expected.Error, actual.Error)
		}
	}
}

func TestAddTool(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	tools := mcpserver.Tools()
	if len(tools) > 0 {
		t.Errorf("expected 0 but got %v", len(tools))
	}

	tool := types.NewTool(
		"tool",
		"tool",
		map[string]any{
			"type":       "object",
			"properties": map[string]any{"path": map[string]string{"type": "string"}},
			"required":   []string{"path"},
		}, nil)

	mcpserver.AddTool(tool)

	tools = mcpserver.Tools()
	if len(tools) != 1 {
		t.Errorf("expected 0 but got %v", len(tools))
	}

	if !reflect.DeepEqual(tools[0], *tool) {
		t.Errorf("expected %v but got %v", *tool, tools[0])
	}
}

func TestAddToolFuncWithInvalidSignature(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	table := []struct {
		handler  any
		expected error
	}{
		{
			map[string]any{},
			ErrHandlerNotFunction,
		},
		{
			func() (*types.ToolResult, error) { return nil, nil },
			ErrHandlerWrongArgs,
		},
		{
			func(_ struct{}) error { return nil },
			ErrHandlerWrongReturns,
		},
		{
			func(_ map[string]any) (*types.ToolResult, error) { return nil, nil },
			ErrHandlerArgNotStruct,
		},
		{
			func(_ struct{}) (*types.ToolResult, any) { return nil, nil },
			ErrHandlerWrongReturns,
		},
		{
			func(_ struct{}) (any, error) { return nil, nil },
			ErrHandlerWrongReturns,
		},
	}

	for _, test := range table {
		err := mcpserver.AddToolFunc("", "", test.handler)

		if err == nil {
			t.Errorf("Expected %#v but got %#v", test.expected, nil)
		}

		if !errors.Is(err, test.expected) {
			t.Errorf("Expected %#v but got %#v", test.expected, err)
		}
	}
}

func TestAddTooFunc(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	expectedInputSpec := map[string]any{
		"type":       "object",
		"properties": map[string]any{"test": map[string]any{"type": "string"}},
		"required":   []string{"test"},
	}

	expectedName := "tool"
	expectedDescription := "tool"

	type ExpectedhandlerArgs struct {
		test string
	}

	expectedResult := types.NewToolResult([]types.OperationContent{*types.NewOperationContent("text", "content", "", nil)})
	expectedHandler := func(arg ExpectedhandlerArgs) (*types.ToolResult, error) {
		return expectedResult, nil
	}

	tools := mcpserver.Tools()
	if len(tools) > 0 {
		t.Errorf("expected 0 but got %v", len(tools))
	}

	mcpserver.AddToolFunc(expectedName, expectedDescription, expectedHandler)

	tools = mcpserver.Tools()
	if len(tools) != 1 {
		t.Errorf("expected 0 but got %v", len(tools))
	}

	if tools[0].Name != expectedName {
		t.Errorf("expected %v but got %v", expectedName, tools[0].Name)
	}

	if tools[0].Description != expectedDescription {
		t.Errorf("expected %v but got %v", expectedDescription, tools[0].Description)
	}

	if !reflect.DeepEqual(tools[0].InputSchema, expectedInputSpec) {
		t.Errorf("expected %v but got %v", expectedInputSpec, tools[0].InputSchema)
	}

	result, err := tools[0].Run(map[string]any{"test": "value"})
	if err != nil {
		t.Errorf("expected not nil but got %v", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("expected %v but got %v", expectedResult, result)
	}
}

func TestAddResource(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	resources := mcpserver.Resources()
	if len(resources) > 0 {
		t.Errorf("expected 0 but got %v", len(resources))
	}

	resource := types.NewResource("Resource", "resource", "uri", nil)

	mcpserver.AddResource(resource)

	resources = mcpserver.Resources()
	if len(resources) != 1 {
		t.Errorf("expected 0 but got %v", len(resources))
	}

	if !reflect.DeepEqual(resources[0], *resource) {
		t.Errorf("expected %v but got %v", *resource, resources[0])
	}
}

func TestHandle(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	table := []struct {
		request          string
		expectedResponse string
	}{
		{
			`{"jsonrpc":"2.0","id":"id1","method":"initialize","params":{"clientName":"testClient","clientVersion":"1.0"}}`,
			fmt.Sprintf(`{"jsonrpc":"2.0","id":"id1","result":{"serverInfo":{"name":"%v","version":"%v"},"capabilities":{"tools":%v,"resources":%v}}}`,
				mcpserver.name, mcpserver.version, len(mcpserver.tools) > 0, len(mcpserver.resources) > 0),
		},
		{
			`{"jsonrpc":"2.0","id":"id2","method":"shutdown","params":{}}`,
			`{"jsonrpc":"2.0","id":"id2","result":{"message":"MCP Session terminated"}}`,
		},
		{
			`{"jsonrpc":"2.0","id":"id3","method":"unknown_method","params":{}}`,
			`{"jsonrpc":"2.0","id":"id3","error":{"code":-32601,"message":"Method Not Found","data":"unknown_method"}}`,
		},
		{
			`{"jsonrpc":"2.0","id":"id4","method":"tools/list","params":{}}`,
			`{"jsonrpc":"2.0","id":"id4","result":{"tools":[]}}`,
		},
		{
			`{"jsonrpc":"2.0","id":"id5","method":"resources/list","params":{}}`,
			`{"jsonrpc":"2.0","id":"id5","result":{"resources":[]}}`,
		},
		{
			`{"jsonrpc":"2.0","id":"id6","method":"tools/call","params":{"name":"non_existing_tool","arguments":{}}}`,
			`{"jsonrpc":"2.0","id":"id6","error":{"code":-32601,"message":"Unknown Tool","data":"tools/call"}}`,
		},
		{
			`{"jsonrpc":"2.0","id":"id7","method":"resources/read","params":{"uri":"non_existing_resource"}}`,
			`{"jsonrpc":"2.0","id":"id7","error":{"code":-32601,"message":"Unknown Resource","data":"resources/read"}}`,
		},
		{
			`invalid_json`,
			`{"jsonrpc":"2.0","id":"","error":{"code":-32700,"message":"Parse error","data":"invalid character 'i' looking for beginning of value"}}`,
		},
	}

	for _, test := range table {
		actualResponse, _ := mcpserver.Handle(test.request)
		if actualResponse != test.expectedResponse {
			t.Errorf("Expected %s but got %s", test.expectedResponse, actualResponse)
		}
	}
}

func TestHandleWithTools(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	expectedInputSpec := map[string]any{
		"type":       "object",
		"properties": map[string]any{"test": map[string]any{"type": "string"}},
		"required":   []string{"test"},
	}

	expectedInputSchema, _ := json.Marshal(expectedInputSpec)
	expectedName := "tool"
	expectedDescription := "tool"

	type ExpectedhandlerArgs struct {
		test string
	}

	expectedResult := types.NewToolResult([]types.OperationContent{*types.NewOperationContent("text", "content", "", nil)})
	expectedHandler := func(arg ExpectedhandlerArgs) (*types.ToolResult, error) {
		return expectedResult, nil
	}

	mcpserver.AddToolFunc(expectedName, expectedDescription, expectedHandler)

	table := []struct {
		request          string
		expectedResponse string
	}{
		{
			`{"jsonrpc":"2.0","id":"id1","method":"initialize","params":{"clientName":"testClient","clientVersion":"1.0"}}`,
			fmt.Sprintf(`{"jsonrpc":"2.0","id":"id1","result":{"serverInfo":{"name":"%v","version":"%v"},"capabilities":{"tools":%v,"resources":%v}}}`,
				mcpserver.name, mcpserver.version, true, false),
		},
		{
			`{"jsonrpc":"2.0","id":"id4","method":"tools/list","params":{}}`,
			fmt.Sprintf(`{"jsonrpc":"2.0","id":"id4","result":{"tools":[{"name":"%v","description":"%v","inputSchema":%v}]}}`,
				expectedName, expectedDescription, string(expectedInputSchema)),
		},
		{
			`{"jsonrpc":"2.0","id":"id6","method":"tools/call","params":{"name":"non_existing_tool","arguments":{}}}`,
			`{"jsonrpc":"2.0","id":"id6","error":{"code":-32601,"message":"Unknown Tool","data":"tools/call"}}`,
		},
		{
			`{"jsonrpc":"2.0","id":"id7","method":"tools/call","params":{"name":"tool","arguments":{"test":"value"}}}`,
			`{"jsonrpc":"2.0","id":"id7","result":{"content":[{"type":"text","text":"content"}]}}`,
		},
	}

	for _, test := range table {
		actualResponse, _ := mcpserver.Handle(test.request)
		if actualResponse != test.expectedResponse {
			t.Errorf("Expected %s but got %s", test.expectedResponse, actualResponse)
		}
	}
}
