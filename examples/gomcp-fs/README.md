# GoMCP File System Tools

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/your-repo)](https://goreportcard.com/report/github.com/your-username/your-repo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue)](https://golang.org/doc/go1.24)

A lightweight **Model Context Protocol (MCP)** server implemented in Go.  
This project exposes three basic file system tools through the MCP interface:

- `ls` — list files in a directory  
- `cd` — change the current working directory  
- `pwd` — print the current working directory  

It uses the [gomcp](https://github.com/mcpunzo/gomcp) library to provide an extensible, protocol-compliant MCP server running over standard input/output.



## 📦 Features

| Command | Description | Parameters |
|----------|--------------|-------------|
| **ls**  | Lists files in the specified directory | `path` *(optional)* — directory path to list (defaults to `"."`) |
| **cd**  | Changes the current working directory | `path` *(required)* — target directory |
| **pwd** | Prints the current working directory | None |



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
cd gomcp/examples/gomcp-fs
```

### 2. Build the binary

```bash
make build
```


## ▶️ Usage

Run the MCP server:

```bash
cd gomcp/examples/gomcp-fs/bin
./gomcp-fs
```

The server will start and wait for MCP-compliant JSON requests via **standard I/O**.  
Example interactions:


## 🧪 Testing

To test locally you can simulate MCP requests by piping JSON messages to stdin or by integrating the server into an MCP-compatible client.

### Initialize

```bash
> echo '{"jsonrpc":"2.0","id":"id1","method":"initialize","params":{"clientName":"testClient","clientVersion":"1.0"}' | ./bin/gomcp-fs
> Starting MCP Server...
> Handling request: initialize
{"jsonrpc":"2.0","id":"id1","result":{"serverInfo":{"name":"gomcp-fs","version":"v1.0.0"},"capabilities":{"tools":true,"resources":false}}}
```

### List Tools

```bash
> echo '{"jsonrpc":"2.0","id":"id4","method":"tools/list","params":{}}' | ./bin/gomcp-fs
> Starting MCP Server...
> Handling request: tools/list
{"jsonrpc":"2.0","id":"id4","result":{"tools":[{"name":"ls","description":"list information about FILEs","inputSchema":{"properties":{"Path":{"type":"string"}},"required":["Path"],"type":"object"}},{"name":"cd","description":"change the current directory","inputSchema":{"properties":{"Path":{"type":"string"}},"required":["Path"],"type":"object"}},{"name":"pwd","description":"print the current working directory","inputSchema":{"properties":{},"required":[],"type":"object"}}]}}
```

### Call tool: ls

```bash
> echo '{"jsonrpc":"2.0","id":"id6","method":"tools/call","params":{"name":"ls"}}' | ./bin/gomcp-fs
> Starting MCP Server...
> Handling request: tools/call
{"jsonrpc":"2.0","id":"id6","result":{"content":[{"type":"text","text":"Makefile"},{"type":"text","text":"README.md"},{"type":"text","text":"bin"},{"type":"text","text":"cmd"}]}}

> echo '{"jsonrpc":"2.0","id":"id6","method":"tools/call","params":{"name":"ls","arguments":{"path":"cmd/"}}}' | ./bin/gomcp-fs
> Starting MCP Server...
> Handling request: tools/call
{"jsonrpc":"2.0","id":"id6","result":{"content":[{"type":"text","text":"Makefile"},{"type":"text","text":"README.md"},{"jsonrpc":"2.0","id":"id6","result":{"content":[{"type":"text","text":"main.go"}]}}
```

### Call tool: cd

```bash
> echo '{"jsonrpc":"2.0","id":"id6","method":"tools/call","params":{"name":"cd","arguments":{"path":"cmd/"}}}' | ./bin/gomcp-fs
> Starting MCP Server...
> Handling request: shutdown
{"jsonrpc":"2.0","id":"id6","result":{"content":[{"type":"text","text":"****/gomcp/examples/gomcp-fs/cmd"}]}}
```

### Call tool: pwd

```bash
> echo '{"jsonrpc":"2.0","id":"id6","method":"tools/call","params":{"name":"pwd"}}' | ./bin/gomcp-fs
> Starting MCP Server...
> Handling request: tools/call
{"jsonrpc":"2.0","id":"id6","result":{"content":[{"type":"text","text":"****/gomcp/examples/gomcp-fs"}]}}
```

### Shutdown

```bash
> echo '{"jsonrpc":"2.0","id":"id2","method":"shutdown","params":{}}' | ./bin/gomcp-fs
> Starting MCP Server...
> Handling request: shutdown
{"jsonrpc":"2.0","id":"id2","result":{"message":"MCP Session terminated"}}
```


## 🛠 Dependencies

- [github.com/mcpunzo/gomcp](https://github.com/mcpunzo/gomcp)

Install manually if needed:

```bash
go get github.com/mcpunzo/gomcp@latest
```


## 👤 Author

**mcpunzo**  
📧 [mcpunzo@gmail.com]  


## 💡 Notes & Recommendations

- Follow [Model Context Protocol](https://modelcontextprotocol.io) conventions for full client compatibility.  
<br>

> **Tip:** This project can serve as a foundation for building custom MCP servers that expose local tools or contextual operations to LLM clients or orchestrators.
