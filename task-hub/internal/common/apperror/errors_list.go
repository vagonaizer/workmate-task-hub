package apperror

// ==================
// Repository errors
// ==================
var (
	ErrRepoNotFound     = New("REPO_NOT_FOUND", "resource not found in repository")
	ErrRepoSaveFailed   = New("REPO_SAVE_FAILED", "failed to save resource in repository")
	ErrRepoDeleteFailed = New("REPO_DELETE_FAILED", "failed to delete resource in repository")
)

// ==================
// Service errors
// ==================
var (
	ErrServiceValidation = New("SERVICE_VALIDATION", "service validation failed")
	ErrServiceConflict   = New("SERVICE_CONFLICT", "resource conflict in service")
)

// ==================
// Transport (HTTP) errors
// ==================
var (
	ErrTransportBadRequest   = New("TRANSPORT_BAD_REQUEST", "bad request")
	ErrTransportUnauthorized = New("TRANSPORT_UNAUTHORIZED", "unauthorized")
	ErrTransportForbidden    = New("TRANSPORT_FORBIDDEN", "forbidden")
	ErrTransportNotFound     = New("TRANSPORT_NOT_FOUND", "resource not found")
	ErrTransportInternal     = New("TRANSPORT_INTERNAL", "internal server error")
)

// ==================
// Application errors
// ==================
var (
	ErrAppInternal = New("APP_INTERNAL", "internal application error")
)
