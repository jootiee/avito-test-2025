package dto

// APIError represents the error response shape per OpenAPI spec
type APIError struct {
	Error ErrorDetail `json:"error"`
}

// ErrorResponse is an alias for APIError for Swagger documentation
type ErrorResponse = APIError

// ErrorDetail contains the error code and message
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewAPIError creates a new API error response
func NewAPIError(code, message string) APIError {
	return APIError{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// Error code constants per OpenAPI spec
const (
	ErrCodeTeamExists        = "TEAM_EXISTS"
	ErrCodePullRequestExists = "PR_EXISTS"
	ErrCodePullRequestMerged = "PR_MERGED"
	ErrCodeNotAssigned       = "NOT_ASSIGNED"
	ErrCodeNoCandidate       = "NO_CANDIDATE"
	ErrCodeNotFound          = "NOT_FOUND"
)
