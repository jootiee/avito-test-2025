package domain

// APIError represents the error response shape per OpenAPI spec
type APIError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// NewAPIError creates a new API error response
func NewAPIError(code, message string) APIError {
	var e APIError
	e.Error.Code = code
	e.Error.Message = message
	return e
}

// Error code constants
const (
	ErrCodeTeamExists  = "TEAM_EXISTS"
	ErrCodePRExists    = "PR_EXISTS"
	ErrCodePRMerged    = "PR_MERGED"
	ErrCodeNotAssigned = "NOT_ASSIGNED"
	ErrCodeNoCandidate = "NO_CANDIDATE"
	ErrCodeNotFound    = "NOT_FOUND"
)
