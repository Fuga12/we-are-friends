package httpapi

// ErrorResponse is the common error envelope for every API endpoint.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError is deliberately small: the frontend can branch on Code while
// Message remains useful to people and logs.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}
