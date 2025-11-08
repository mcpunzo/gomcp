# GoMCP Simple Calculator

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/your-repo)](https://goreportcard.com/report/github.com/your-username/your-repo)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue)](https://golang.org/doc/go1.24)

A simple **Model Context Protocol (MCP)** server built in Go, demonstrating how to expose basic arithmetic operations (`plus` and `minus`) as callable tools using the [gomcp](https://github.com/mcpunzo/gomcp) framework.


## 🚀 Overview

`gomcp-calculator` showcases how to create an MCP-compliant server using **Go** and **gomcp**, leveraging the HTTP transport layer to expose tools that perform basic calculations.

Each tool can be invoked via JSON-RPC calls over HTTP, making it easy to integrate this server with MCP clients or other systems that implement the MCP specification.


## 📦 Features

- ✅ **Implements MCP protocol** using `gomcp`
- 🌐 **HTTP transport layer** (default port: `8080`)
- ➕ **Addition** and ➖ **Subtraction** tools
- 🔍 **Structured JSON-RPC communication**
- ⚙️ **Easy extensibility** for adding new tools




## 🚀 Getting Started

### Prerequisites

- **Go 1.24.1+**
- A compatible **MCP client** that communicates via `stdin/stdout` transport.

Check your Go installation:

```bash
go version
```

### 1. Clone the repository

```bash
git clone git@github.com:mcpunzo/gomcp.git
cd gomcp/examples/gomcp-calculator
```

### 2. Build the binary

```bash
make build
```


## ▶️ Usage

Run the MCP server:

```bash
cd gomcp/examples/gomcp-calculator/bin
./gomcp-calculator
```

The server will start and wait for MCP-compliant JSON requests via **HTTP**.  


## 🧪 Testing

To test locally you can simulate MCP requests by sending JSON messages via HTTP.

### Prerequisites

Start the server as described above.
In another shell:

### Initialize

```bash
> curl -X POST http://localhost:8080/mcp -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":"id1","method":"initialize","params":{"clientName":"testClient","clientVersion":"1.0"}}'
{"jsonrpc":"2.0","id":"id1","result":{"serverInfo":{"name":"gomcp-calculator","version":"v1.0.0"},"capabilities":{"tools":true,"resources":false}}}
```

### List Tools

```bash
> curl -X POST http://localhost:8080/mcp -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":"id4","method":"tools/list","params":{}}'                                         
{"jsonrpc":"2.0","id":"id4","result":{"tools":[{"name":"plus","description":"Sum operator for 2 int parameters","inputSchema":{"properties":{"A":{"type":"int"},"B":{"type":"int"}},"required":["A","B"],"type":"object"}},{"name":"minus","description":"Minus operator for 2 int parameters","inputSchema":{"properties":{"A":{"type":"int"},"B":{"type":"int"}},"required":["A","B"],"type":"object"}}]}}
```

### Call tool: plus

```bash
> curl -X POST http://localhost:8080/mcp -H "Content-Type: application/json" -d '{"jsonrpc": "2.0", "id": "1", "method": "tools/call", "params": {"name": "plus", "arguments": { "a": 5, "b": 3 }}}'
{"jsonrpc":"2.0","id":"1","result":{"content":[{"type":"text","text":"8"}]}}
```

### Call tool: minus

```bash
> curl -X POST http://localhost:8080/mcp -H "Content-Type: application/json" -d '{"jsonrpc": "2.0", "id": "1", "method": "tools/call", "params": {"name": "minus", "arguments": { "a": 5, "b": 3 }}}'
{"jsonrpc":"2.0","id":"1","result":{"content":[{"type":"text","text":"2"}]}}
```

### Shutdown

```bash
> curl -X POST http://localhost:8080/mcp -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":"id2","method":"shutdown","params":{}}'
{"jsonrpc":"2.0","id":"id2","result":{"message":"MCP Session terminated"}}
```

## 🧩 Extending the server
You can easily extend the MCP server by adding new tools following the same pattern:
```go
func addMultiplyTool(mcp *gomcp.MCPServer) {
	mcp.AddToolFunc("multiply", "Multiplication of two integers", func(params CalculatorParams) (*types.ToolResult, error) {
		content := []types.OperationContent{
			*types.NewOperationContent("text", strconv.Itoa(params.A*params.B), "", nil),
		}
		return types.NewToolResult(content), nil
	})
}
```
Then, just call:

```go
addMultiplyTool(mcp)
```

## 👤 Author

**mcpunzo**  
📧 [mcpunzo@gmail.com]  


## 💡 Notes & Recommendations

- Follow [Model Context Protocol](https://modelcontextprotocol.io) conventions for full client compatibility.  
<br>

> **Tip:** This project can serve as a foundation for building custom MCP servers that expose local tools or contextual operations to LLM clients or orchestrators.
