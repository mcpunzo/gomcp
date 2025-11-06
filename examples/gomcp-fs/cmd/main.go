package main

import (
	"errors"
	"log"
	"os"

	"github.com/mcpunzo/gomcp"
	"github.com/mcpunzo/gomcp/internal/transport"
	"github.com/mcpunzo/gomcp/types"
)

type FSReaderParams struct {
	Path string `json:"path"`
}

func main() {
	mcp := gomcp.New("gomcp-fs", "v1.0.0").WithTransport(transport.NewStdIOTransport())
	addLsTool(mcp)
	addCdTool(mcp)
	addPwdTool(mcp)

	mcp.Run()
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
