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

	ErrDeleteFailed = ErrorMessage{
		Message: "Failed to delete data",
		Details: []string{
			"The server failed to delete the requested resource from the database.",
			"This might be due to a database error or the resource may not exist.",
			"Please try again or contact support if the issue persists.",
		},
	}

	ErrEmptyContent = ErrorMessage{
		Message: "Empty content field",
		Details: []string{
			"The content field must not be empty.",
			"Make sure to fill in the required fields in the request body.",
			"Try again after entering valid content.",
		},
	}

	ErrCommentSave = ErrorMessage{
		Message: "Failed to save comment",
		Details: []string{
			"The comment could not be saved to the database.",
			"Check for invalid data or internal server error.",
			"If the problem persists, please contact support.",
		},
	}

	ErrUnauthorized = ErrorMessage{
		Message: "Unauthorized access",
		Details: []string{
			"You are not authorized to access or modify this resource.",
			"Make sure you have the correct permissions.",
			"This action can only be performed by the assigned lecturer or admin.",
		},
	}

	ErrNoComments = ErrorMessage{
		Message: "No comments found",
		Details: []string{
			"There are no comments for this assignment.",
			"Comments may not have been submitted yet.",
			"Try again later after someone has added a comment.",
		},
	}

	ErrDeleteComment = ErrorMessage{
		Message: "Failed to delete comment(s)",
		Details: []string{
			"The system failed to delete comment(s) associated with the assignment.",
			"This may be due to a database error or missing comment(s).",
			"Please check the assignment ID or try again later.",
		},
	}

	ErrCommentNotFound = ErrorMessage{
		Message: "Comment(s) not found",
		Details: []string{
			"No comment(s) associated with the specified assignment ID were found.",
			"Please check if the assignment has any comments.",
			"Ensure the assignment ID is correct.",
		},
	}

	ErrAssignmentNotUpdated = ErrorMessage{
		Message: "Assignment not updated",
		Details: []string{
			"No changes were made to the assignment.",
			"The submitted data might be the same as existing values.",
			"Ensure the assignment ID is valid and the data is different.",
		},
	}

	ErrSubmissionStatusNotFound = ErrorMessage{
		Message: "Submission status not found",
		Details: []string{
			"The specified submission status ID could not be found.",
			"This might indicate an invalid status ID provided or a data integrity issue.",
			"Please ensure the submission status is valid. If the problem persists, contact support.",
		},
	}

	ErrAssignmentDeadlinePassed = ErrorMessage{
		Message: "Assignment deadline has passed",
		Details: []string{
			"The submission could not be processed because the assignment deadline has already passed.",
			"Please check the assignment details for the correct submission period.",
			"Late submissions are not allowed for this assignment.",
		},
	}

	ErrSubmissionFailed = ErrorMessage{
		Message: "Failed to submit assignment",
		Details: []string{
			"There was an issue processing your assignment submission.",
			"This might be due to a problem with your submission data or a temporary server issue.",
			"Please check your input and try again. If the problem persists, contact support.",
		},
	}

	ErrSubmissionAlreadyExists = ErrorMessage{
		Message: "Submission already exists",
		Details: []string{
			"A submission for this assignment by this student already exists.",
			"Only one submission is allowed per student for this assignment.",
			"If you need to update your submission, please use the update endpoint or contact the lecturer.",
		},
	}
)
