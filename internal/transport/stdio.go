package transport

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/mcpunzo/gomcp"
)

type StdioTransport struct {
	mgp *gomcp.MCPServer
}

func NewStdIOTransport() *StdioTransport {
	return &StdioTransport{}
}

// SetMCPServer sets the MCPServer for the StdioTransport.
func (s *StdioTransport) SetMCPServer(mcpserver *gomcp.MCPServer) {
	s.mgp = mcpserver
}

// Start starts the StdioTransport to read from stdin and write to stdout.
func (s *StdioTransport) Start() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Errore: %v\n", err)
			break
		}

		response, _ := s.mgp.Handle(line)
		fmt.Fprintln(writer, response)
		writer.Flush()

	}
}
