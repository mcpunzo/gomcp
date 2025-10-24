package types

type ShutdownParams struct{}

// ShutdownResult represents the result of a shutdown operation.
type ShutdownResult struct {
	// A message indicating the shutdown status.
	Message string `json:"message"`
}

// NewShutdownParams creates a new ShutdownParams.
func NewShutdownParams() *ShutdownParams {
	return &ShutdownParams{}
}

// NewShutdownResult creates a new ShutdownResult with the given message.
func NewShutdownResult(message string) *ShutdownResult {
	return &ShutdownResult{Message: message}
}
