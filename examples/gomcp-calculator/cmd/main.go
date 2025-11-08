package main

import (
	"log"
	"strconv"

	"github.com/mcpunzo/gomcp"
	"github.com/mcpunzo/gomcp/internal/transport"
	"github.com/mcpunzo/gomcp/types"
)

type CalculatorParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	mcp := gomcp.New("gomcp-calculator", "v1.0.0").WithTransport(transport.NewHttpTransport(8080))
	addPlusTool(mcp)
	addMinusTool(mcp)
	mcp.Run()
}

func addPlusTool(mcp *gomcp.MCPServer) {
	mcp.AddToolFunc("plus", "Sum operator for 2 int parameters", func(params CalculatorParams) (*types.ToolResult, error) {
		log.Print(params)

		content := []types.OperationContent{*types.NewOperationContent("text", strconv.Itoa(params.A+params.B), "", nil)}

		return types.NewToolResult(content), nil
	})
}

func addMinusTool(mcp *gomcp.MCPServer) {
	mcp.AddToolFunc("minus", "Minus operator for 2 int parameters", func(params CalculatorParams) (*types.ToolResult, error) {
		log.Print(params)

		content := []types.OperationContent{*types.NewOperationContent("text", strconv.Itoa(params.A-params.B), "", nil)}

		return types.NewToolResult(content), nil
	})
}
