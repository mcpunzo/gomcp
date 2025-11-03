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
	addCdTool(mcp)
	addPwdTool(mcp)

	handleRequests(mcp, types.NewJSONRPCRequest("id1", gomcp.Initialize, types.NewInitializeParams("test", "v0.0.1")))
	handleRequests(mcp, types.NewJSONRPCRequest("id2", gomcp.ListResources, nil))
	handleRequests(mcp, types.NewJSONRPCRequest("id3", gomcp.ListTools, nil))

	handleRequests(mcp, types.NewJSONRPCRequest("pwd1",
		gomcp.CallTool,
		types.NewCallToolParams("pwd", nil)))

	handleRequests(mcp, types.NewJSONRPCRequest("ls1",
		gomcp.CallTool,
		types.NewCallToolParams("ls", map[string]any{"Path": "../"})))

	handleRequests(mcp, types.NewJSONRPCRequest("cd1",
		gomcp.CallTool,
		types.NewCallToolParams("cd", map[string]any{"Path": "../"})))

	handleRequests(mcp, types.NewJSONRPCRequest("pwd2",
		gomcp.CallTool,
		types.NewCallToolParams("pwd", nil)))

	handleRequests(mcp, types.NewJSONRPCRequest("finish", gomcp.Shutdown, nil))
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

func addCdTool(mcp *gomcp.MCPServer) {
	cd_handler := func(params FSReaderParams) (*types.ToolResult, error) {
		log.Print(params.Path)

		err := os.Chdir(params.Path)
		if err != nil {
			return nil, err
		}

		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		content := []types.OperationContent{
			*types.NewOperationContent("text", currentDir, "", nil),
		}

		return types.NewToolResult(content), nil
	}

	mcp.AddToolFunc("cd", "change the current directory", cd_handler)
}

func addPwdTool(mcp *gomcp.MCPServer) {
	pwd_handler := func(_ struct{}) (*types.ToolResult, error) {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		content := []types.OperationContent{
			*types.NewOperationContent("text", currentDir, "", nil),
		}

		return types.NewToolResult(content), nil
	}

	mcp.AddToolFunc("pwd", "print the current working directory", pwd_handler)
}
