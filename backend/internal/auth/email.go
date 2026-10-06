package auth

const (
	CredentialEmail  = "email"
	CredentialGoogle = "google"
	CredentialGithub = "github"
)

// 60 days
const MaxRefreshCookieMaxAge = 60 * 24 * 60 * 60

// The shape of the JSON body expected when login by email.
type EmailLoginBody struct {
	Email string `json:"email" binding:"required,email"`
}

// The shape of JSON body expected when verifying login by email
type EmailVerifyBody struct {
	Token string `json:"token" binding:"required"`
}
