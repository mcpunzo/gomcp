package types

// OperationContent represents the content returned by resource operations.
type OperationContent struct {
	Type string `json:"type"`           // text, markdown, image, json, uri, ecc.
	Text string `json:"text,omitempty"` // "text"/"markdown"
	Data any    `json:"data,omitempty"` // JSON or binari
	URI  string `json:"uri,omitempty"`  // external references
}
