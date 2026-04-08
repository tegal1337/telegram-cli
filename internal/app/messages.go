package app

// FocusChangedMsg signals a focus panel change.
type FocusChangedMsg struct {
	Panel FocusPanel
}

// ErrorMsg carries an error to display.
type ErrorMsg struct {
	Err error
}

// AuthStateChangedMsg is sent from the authorizer callback.
type AuthStateChangedMsg struct {
	State int
	Hint  string
}

// AuthenticatedMsg signals that authentication is complete.
type AuthenticatedMsg struct {
	UserId    int64
	FirstName string
	LastName  string
}
