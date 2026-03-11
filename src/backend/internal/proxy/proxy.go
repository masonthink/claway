package proxy

// LLMProxy handles forwarding requests to LLM providers
// while tracking token usage and enforcing limits.
type LLMProxy struct {
	// TODO: add provider clients, rate limiters, etc.
}

// New creates a new LLMProxy instance.
func New() *LLMProxy {
	return &LLMProxy{}
}
