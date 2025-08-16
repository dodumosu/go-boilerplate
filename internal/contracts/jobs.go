package contracts

const (
	TypeVerificationEmail = "email:verify"
)

type VerificationEmailPayload struct {
	Email            string
	Name             string
	VerificationLink string
	CorrelationID    string
}
