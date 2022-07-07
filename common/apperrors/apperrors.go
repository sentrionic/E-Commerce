package apperrors

// Account Errors
const (
	InvalidOldPassword  = "Invalid old password"
	InvalidCredentials  = "Invalid email and password combination"
	DuplicateEmail      = "An account with that email already exists"
	PasswordsDoNotMatch = "Passwords do not match"
	InvalidResetToken   = "Invalid reset token"
)

// Generic Errors
const (
	InvalidSession = "Provided session is invalid"
	ServerError    = "Something went wrong. Try again later"
	Unauthorized   = "Not Authorized"
)
