package composer

// MessageSubmittedMsg is emitted when the user submits a message.
type MessageSubmittedMsg struct {
	ChatId         int64
	Text           string
	ReplyToId      int64
	EditMessageId  int64
}

// AttachmentAddedMsg is emitted when a file is attached.
type AttachmentAddedMsg struct {
	FilePath string
	FileType string // "photo", "document", "video", "voice"
}
