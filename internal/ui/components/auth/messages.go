package auth

// AuthCompleteMsg is emitted when authentication succeeds.
type AuthCompleteMsg struct{}

// PhoneSubmittedMsg is emitted when the user submits their phone number.
type PhoneSubmittedMsg struct {
	Phone string
}

// CodeSubmittedMsg is emitted when the user submits the verification code.
type CodeSubmittedMsg struct {
	Code string
}

// PasswordSubmittedMsg is emitted when the user submits their 2FA password.
type PasswordSubmittedMsg struct {
	Password string
}
