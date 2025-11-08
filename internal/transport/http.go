package transport

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mcpunzo/gomcp"
)

type HttpTransport struct {
	mgp  *gomcp.MCPServer
	port int
}

func NewHttpTransport(port int) *HttpTransport {
	return &HttpTransport{port: port}
}

// SetMCPServer sets the MCPServer for the StdioTransport.
func (h *HttpTransport) SetMCPServer(mcpserver *gomcp.MCPServer) {
	h.mgp = mcpserver
}

// Start starts the HTTP server to read from a post to /mcp endpoint.
func (h *HttpTransport) Start() {
	http.HandleFunc("/mcp", h.handler)

	port := fmt.Sprintf(":%d", h.port)
	log.Printf("Server started and listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (h *HttpTransport) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading the request body", http.StatusBadRequest)
		return
	}

	bodyString := string(bodyBytes)

	response, err := h.mgp.Handle(bodyString)
	log.Printf("Response: %s", response)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, response)

}
