# GoMCP — Model Context Protocol Framework for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/mcpunzo/gomcp)](https://goreportcard.com/report/github.com/mcpunzo/gomcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue)](https://golang.org/doc/go1.24)

**GoMCP** is a lightweight, extensible framework for building **Model Context Protocol (MCP)** servers in Go.  
It provides a simple way to define, register, and expose tools and resources over the MCP protocol — a JSON-RPC–based interface designed for seamless communication between AI models and context providers.



## 🧩 Key Features

- ⚙️ **JSON-RPC compliant** message structure (requests, responses, and errors)  
- 🔌 **Transport abstraction layer** — easily switch between `stdio`, TCP, or custom transports  
- 🧠 **Dynamic tool registration** with automatic JSON schema generation  
- 🧾 **Resources API** to expose structured contextual data  
- 🧰 **Extensible architecture** suitable for building custom MCP servers or model integrations  


## 🚀 Installation

Make sure you have **Go 1.24.1+** installed:

```bash
go version
```

Then install the package:

```bash
go get github.com/mcpunzo/gomcp@latest
```

## 🧠 Framework Overview

### MCPServer

The core of GoMCP is the `MCPServer` struct. It manages:
- Tool and resource registration  
- JSON-RPC message parsing and dispatching  
- Error handling and protocol compliance  

Example initialization:
```go
mcp := gomcp.New("my-server", "v1.0.0").WithTransport(transport.NewStdIOTransport())
mcp.Run()
```

### Adding Tools

Each tool is registered using `AddToolFunc`:

```go
mcp.AddToolFunc("echo", "Echoes a message", func(params struct{ Message string `json:"message"` }) (*types.ToolResult, error) {
    content := []types.OperationContent{
        *types.NewOperationContent("text", params.Message, "", nil),
    }
    return types.NewToolResult(content), nil
})
```

The framework automatically validates the function signature, converts JSON arguments into the provided struct, and generates a JSON schema for the tool.

### Built-in JSON-RPC Methods

| Method | Description |
|---------|-------------|
| `initialize` | Starts an MCP session |
| `shutdown` | Gracefully terminates the session |
| `tools/list` | Lists all registered tools |
| `tools/call` | Invokes a tool with the provided parameters |
| `resources/list` | Lists available resources |
| `resources/read` | Reads a specific resource by URI |


## ⚙️ Architecture

GoMCP is composed of several layers:

| Layer | Description |
|--------|--------------|
| **Transport** | Defines communication (e.g. `StdIOTransport`, TCP, WebSocket). |
| **Server Core** | Parses JSON-RPC messages and dispatches requests. |
| **Tools & Resources** | Domain-specific capabilities registered dynamically. |
| **Types Package** | Contains JSON-RPC and MCP data structures. |


## 📘 Quick Start Example

Below is a minimal MCP server exposing the 'ls' filesystem tools via STDIO.
You can find more examples in the [exemple directory](./examples).


```go
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
    mcp.Run()
}

func addLsTool(mcp *gomcp.MCPServer) {
    ls_handler := func(params FSReaderParams) (*types.ToolResult, error) {
        log.Print(params.Path)

        path := params.Path
        if path == "" {
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

```

## 🧰 Adding New Tools

To add your own tool:

1. Define the parameter struct  
2. Implement a handler that returns `(*types.ToolResult, error)`  
3. Register it via `AddToolFunc(name, description, handler)`  

Example:

```go
func addHelloTool(mcp *gomcp.MCPServer) {
    handler := func(params struct{ Name string `json:"name"` }) (*types.ToolResult, error) {
        message := fmt.Sprintf("Hello, %s!", params.Name)
        content := []types.OperationContent{
            *types.NewOperationContent("text", message, "", nil),
        }
        return types.NewToolResult(content), nil
    }

    mcp.AddToolFunc("hello", "Greet a user by name", handler)
}
```

## 🛠 Error Handling

GoMCP provides standardized error codes (based on JSON-RPC conventions):

| Code | Constant | Meaning |
|------|-----------|----------|
| `-32700` | `ErrParse` | Invalid JSON received |
| `-32600` | `ErrInvalidRequest` | Invalid request object |
| `-32601` | `ErrMethodNotFound` | Method not found |
| `-32602` | `ErrInvalidParams` | Invalid parameters |
| `-32603` | `ErrInternal` | Internal server error |
| `-32000+` | `ErrServerGeneric`, `ErrAccessDenied`, `ErrNotFound` | Custom server errors |


## 📦 Dependencies

- [Go Standard Library](https://pkg.go.dev/std)
- [github.com/mcpunzo/gomcp](https://github.com/mcpunzo/gomcp)


## 📄 License

This project is licensed under the **MIT License**.  
See the [LICENSE](./LICENSE) file for details.


## 👤 Author

**M.C. Punzo**  
📧 [mcpunzo@gmail.com]  
🔗 [@mcpunzo](https://github.com/mcpunzo)


## 💡 Notes

- GoMCP follows the [Model Context Protocol](https://modelcontextprotocol.io) specification.  
- The design encourages modularity — each tool can be added as an independent Go function.  
- This framework can serve as a foundation for creating MCP-compatible agents, extensions, or AI backends.


> **Tip:** GoMCP is ideal for integrating structured Go logic with LLMs or context-aware AI systems using standardized MCP communication.
