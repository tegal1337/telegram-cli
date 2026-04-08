package composer

// MessageSubmittedMsg is emitted when the user submits a message.
type MessageSubmittedMsg struct {
	ChatID         int64
	Text           string
	ReplyToID      int64
	EditMessageID  int64
}

// AttachmentAddedMsg is emitted when a file is attached.
type AttachmentAddedMsg struct {
	FilePath string
	FileType string // "photo", "document", "video", "voice"
}
