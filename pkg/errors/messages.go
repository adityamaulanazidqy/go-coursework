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

	ErrAssignmentNotFound = ErrorMessage{
		Message: "Assignment Not Found",
		Details: []string{
			"Assignment ID is required in the request URL.",
			"Make sure the endpoint follows the correct format, e.g., /assignments/{id}.",
			"Check if the assignment ID is missing or malformed.",
		},
	}

	ErrBodyParse = ErrorMessage{
		Message: "Failed to parse body",
		Details: []string{
			"One of the requests is not eligible",
		},
	}

	ErrFileUpload = ErrorMessage{
		Message: "Failed to upload file",
		Details: []string{
			"The file could not be uploaded or was not found in the request.",
			"Make sure to include a valid file in the 'file' form field.",
			"Supported formats and size limits may apply.",
		},
	}

	ErrFileOpen = ErrorMessage{
		Message: "Failed to open uploaded file",
		Details: []string{
			"The uploaded file could not be opened for processing.",
			"This might be due to file corruption or internal server error.",
			"Please try uploading the file again.",
		},
	}

	ErrFileSave = ErrorMessage{
		Message: "Failed to save file",
		Details: []string{
			"The uploaded file could not be saved to the server.",
			"This may be due to permission issues, unsupported format, or internal errors.",
			"Please try again or contact support if the problem persists.",
		},
	}
)
