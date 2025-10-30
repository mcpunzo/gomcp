package gomcp

import (
	"errors"
	"log"
	"reflect"
	"strings"
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
			types.NewJSONRPCResponse("id", types.NewShutdownResult(mcpserver.shutdownMessage), nil),
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
		func(args any) (*types.ToolResult, error) {
			params, ok := args.(map[string]any)
			if !ok {
				return nil, err
			}

			path, ok := params["path"].(string)
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

	result, err := tools[0].Run(ExpectedhandlerArgs{test: "value"})
	if err != nil {
		t.Errorf("expected not nil but got %v", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("expected %v but got %v", expectedResult, result)
	}
}

func TestAddToolFunc_Panics(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	tests := []struct {
		name      string
		handler   any
		wantPanic string
	}{
		{
			name:      "non-function handler",
			handler:   123,
			wantPanic: "handler must be a function",
		},
		{
			name: "no parameter handler",
			handler: func() (*types.ToolResult, error) {
				return nil, nil
			},
			wantPanic: "handler need to have a parameter of type struct",
		},
		{
			name: "parameter not struct",
			handler: func(x int) (*types.ToolResult, error) {
				return nil, nil
			},
			wantPanic: "handler parameter need to be of type struct",
		},
		{
			name: "pointer to non-struct",
			handler: func(x *int) (*types.ToolResult, error) {
				return nil, nil
			},
			wantPanic: "handler parameter need to be of type struct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic %q, but did not panic", tt.wantPanic)
				} else {
					rstr := r.(string)
					if !contains(rstr, tt.wantPanic) {
						t.Errorf("expected panic %q, but got %q", tt.wantPanic, rstr)
					}
				}
			}()

			mcpserver.AddToolFunc("test", "desc", tt.handler)
		})
	}
}

func TestAddToolFunc_ValidHandler(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	type MyArgs struct {
		Foo string
		Bar int
	}

	expectedResult := types.NewToolResult([]types.OperationContent{*types.NewOperationContent("text", "content", "", nil)})

	handler := func(a MyArgs) (*types.ToolResult, error) {
		if a.Foo != "foo" || a.Bar != 42 {
			return nil, errors.New("wrong args")
		}
		return expectedResult, nil
	}

	mcpserver.AddToolFunc("tool1", "descrizione", handler)

	tool, ok := mcpserver.tools["tool1"]
	if !ok {
		t.Fatalf("tool not found in map")
	}

	expectedSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"Foo": map[string]any{"type": "string"},
			"Bar": map[string]any{"type": "int"},
		},
		"required": []string{"Foo", "Bar"},
	}

	if !reflect.DeepEqual(tool.InputSchema, expectedSchema) {
		t.Errorf("schema mismatch:\nexpected %#v\ngot %#v", expectedSchema, tool.InputSchema)
	}

	run := reflect.ValueOf(tool.Run)
	if run.Kind() != reflect.Func {
		t.Fatalf("Run is not a function")
	}

	args := MyArgs{Foo: "foo", Bar: 42}
	out, err := tool.Run(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != expectedResult {
		t.Errorf("expected result %v, got %v", expectedResult, out)
	}
}

func TestAddToolFunc_InvalidHandlerReturn(t *testing.T) {
	mcpserver, teardown := setupTest(t)
	defer teardown(t)

	type Args struct{ Foo string }

	tests := []struct {
		name    string
		handler any
	}{
		{
			name:    "no return values",
			handler: func(a Args) {},
		},
		{
			name: "one return value only",
			handler: func(a Args) *types.ToolResult {
				return nil
			},
		},
		{
			name: "wrong return types",
			handler: func(a Args) (string, int) {
				return "oops", 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcpserver.AddToolFunc("badtool", "desc", tt.handler)

			tool := mcpserver.tools["badtool"]
			if tool == nil {
				t.Fatalf("tool not registered")
			}

			// Prova a eseguire il Run: deve ritornare errore
			_, err := tool.Run(Args{Foo: "x"})
			if err == nil {
				t.Errorf("expected error for handler with invalid return types")
			} else if !strings.Contains(err.Error(), "handler must return (*types.ToolResult, error)") {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && reflect.DeepEqual(s[:len(substr)], substr)
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
