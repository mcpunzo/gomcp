package types

// JSONRPCRequest represents a JSON-RPC request object.
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC response object.
type JSONRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Id      string           `json:"id"`
	Result  any              `json:"result,omitempty"`
	Error   *JSONRPCErrorObj `json:"error,omitempty"`
}

// JSONRPCErrorObj represents a JSON-RPC error object.
type JSONRPCErrorObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// NewJSONRPCRequest creates a new JSON-RPC request object.
func NewJSONRPCRequest(id, method string, params any) *JSONRPCRequest {
	return &JSONRPCRequest{JSONRPC: "2.0", Id: id, Method: method, Params: params}
}

// NewJSONRPCResponse creates a new JSON-RPC response object.
func NewJSONRPCResponse(id string, result any, err *JSONRPCErrorObj) *JSONRPCResponse {
	return &JSONRPCResponse{JSONRPC: "2.0", Id: id, Result: result, Error: err}
}

// NewJSONRPCErrorObj creates a new JSON-RPC error object.
func NewJSONRPCErrorObj(code int, msg string, data any) *JSONRPCErrorObj {
	return &JSONRPCErrorObj{Code: code, Message: msg, Data: data}
}
