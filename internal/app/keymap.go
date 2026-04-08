package app

// FocusPanel identifies which UI panel has focus.
type FocusPanel int

const (
	PanelChatList FocusPanel = iota
	PanelChatView
	PanelComposer
	PanelSearch
	PanelContacts
	PanelGroupInfo
)

// ScreenState identifies the current top-level screen.
type ScreenState int

const (
	ScreenAuth ScreenState = iota
	ScreenLoading
	ScreenMain
)
