package chatview

type ScrollToBottomMsg struct{}

type LoadMoreHistoryMsg struct {
	ChatId int64
}

type MessageActionMsg struct {
	Action    string // "reply", "edit", "delete", "forward"
	ChatId    int64
	MessageId int64
}

// MediaPlayMsg is sent when media playback starts.
type MediaPlayMsg struct {
	Status string // "playing", "downloading", "error", "opened"
	Info   string
}
