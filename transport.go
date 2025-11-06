package gomcp

type Transport interface {
	SetMCPServer(mcpserver *MCPServer)
	Start()
}
