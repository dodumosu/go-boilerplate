package contracts

const (
	TypePasswordResetEmail = "email:password-reset"
	TypeVerificationEmail  = "email:verify"
)

type PasswordResetEmailPayload struct {
	Email         string
	Name          string
	ResetLink     string
	CorrelationID string
}

type VerificationEmailPayload struct {
	Email            string
	Name             string
	VerificationLink string
	CorrelationID    string
}
