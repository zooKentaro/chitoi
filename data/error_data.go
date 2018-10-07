package data

// ErrorBadRequest is XXX
type ErrorBadRequest struct{}

func (e *ErrorBadRequest) Error() string {
    return "bad request"
}
