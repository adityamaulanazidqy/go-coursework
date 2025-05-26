package errors

type ErrorMessage struct {
	Message string
	Details []string
}

var (
	ErrMissingClaims = ErrorMessage{
		Message: "Missing claims",
		Details: []string{"There was an error with your token. Please re-login to get a new access token"},
	}

	ErrInternalServer = ErrorMessage{
		Message: "Internal Server Error",
		Details: []string{"User Internal Server Error", "Make sure your internet is on", "There may be a problem with our server"},
	}

	ErrInvalidDeadline = ErrorMessage{
		Message: "Invalid deadline format",
		Details: []string{"Format deadline must '2006-01-02 15:04:05'"},
	}

	ErrUserNotFound = ErrorMessage{
		Message: "User Not Found",
		Details: []string{"Please enter the user ID correctly"},
	}

	ErrMissingAsgnID = ErrorMessage{
		Message: "Missing assignment ID",
		Details: []string{
			"Assignment ID is required in the request URL.",
			"Make sure the endpoint follows the correct format, e.g., /assignments/{id}/comments.",
			"Check if the assignment ID is missing or malformed.",
		},
	}
)
