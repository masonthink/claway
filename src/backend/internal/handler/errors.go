package handler

import "strings"

// userMessage extracts a user-friendly error message from a potentially
// wrapped error chain. Internal details like "get idea by id:" are stripped
// to avoid leaking implementation info in API responses.
func userMessage(err error) string {
	msg := err.Error()
	// Take only the last segment after the final ": " separator,
	// which is typically the root cause message.
	if i := strings.LastIndex(msg, ": "); i != -1 {
		msg = msg[i+2:]
	}
	return msg
}
