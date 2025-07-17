package telegram

// MessageType represents tpye of Telegram message
type MessageType string

const (
	// MessageTypeRegular is 'message'
	MessageTypeRegular = "message"

	// MessageTypeEdited is 'edited_message'
	MessageTypeEdited = "edited_message"

	// MessageTypeChannelPost is 'channel_post'
	MessageTypeChannelPost = "channel_post"

	// MessageTypeEditedChannelPost is 'edited_channel_post'
	MessageTypeEditedChannelPost = "edited_channel_post"
)
