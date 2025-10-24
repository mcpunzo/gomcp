package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mcpunzo/gomcp"
	"github.com/mcpunzo/gomcp/types"
)

type FSReaderParams struct {
	Path string `json:"path"`
}

func main() {
	mcp := gomcp.New("name", "v1.0.0")

	addLsTool(mcp)

	handleRequests(mcp, types.NewJSONRPCRequest("id1", gomcp.Initialize, types.NewInitializeParams("test", "v0.0.1")))
	handleRequests(mcp, types.NewJSONRPCRequest("id2", gomcp.ListResources, nil))
	handleRequests(mcp, types.NewJSONRPCRequest("id3", gomcp.ListTools, nil))

	handleRequests(mcp, types.NewJSONRPCRequest("id",
		gomcp.CallTool,
		types.NewCallToolParams("ls", map[string]any{"Path": "../"})))
}

func handleRequests(mcpserver *gomcp.MCPServer, request *types.JSONRPCRequest) {
	response := mcpserver.HandleRequest(request)

	res, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("res: %+s\n", res)
}

func addLsTool(mcp *gomcp.MCPServer) {
	ls_handler := func(params FSReaderParams) (*types.ToolResult, error) {
		log.Print(params.Path)

		path := params.Path
		if params.Path == "" {
			path = "."
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if !fileInfo.IsDir() {
			return nil, errors.New("the specified path is not a directory")
		}

		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		content := []types.OperationContent{}

		for _, file := range files {
			content = append(content, *types.NewOperationContent("text", file.Name(), "", nil))
		}

		return types.NewToolResult(content), nil
	}

	mcp.AddToolFunc("ls", "list information about FILEs", ls_handler)
}
