package app

// FocusChangedMsg signals a focus panel change.
type FocusChangedMsg struct {
	Panel FocusPanel
}

// ErrorMsg carries an error to display.
type ErrorMsg struct {
	Err error
}

// AuthenticatedMsg signals that authentication is complete.
type AuthenticatedMsg struct {
	UserID    int64
	FirstName string
	LastName  string
}
