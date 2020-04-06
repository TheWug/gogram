package data

import (
)

type TargetData struct {
	ChatId               interface{}      // the chat to send the message to. (int, int64, or string)
}

type SourceData struct {
	SourceChatId    ChatID
	SourceMessageId MsgID
	SourceInlineId  InlineID
}

type SendData struct {
	TargetData
	Text                 string           // for ordinary messages, the text. For media messages, the caption.
	ParseMode           *MessageParseMode // the parse mode of the text/caption for this message, or nil for default behavior
	DisableNotification *MessageNotify    // the notification behavior for this message, or nil for default behavior
	ReplyToId           *MsgID            // the message to reply to, or nil not to reply to another message
	ReplyMarkup          interface{}      // reply markup, or nil for no markup (InlineKeyboard, ReplyKeyboard, ReplyKeyboardRemove, ForceReply)
}

type MediaData struct {
	File                 interface{} // the file to send. (string, FileId, []byte, io.Reader, reqtify.FormFile)
	FileName             string      // an optional filename to include. ignored if a string, FileId, or reqtify.FormFile is used for File.
	Thumb                interface{} // the thumbnail to send. (string, []byte, io.Reader, reqtify.FormFile)
}

type LengthData struct {
	Duration int
}

type ResolutionData struct {
	Width, Height int
}

type ArtistData struct {
	Performer, Title string
}

type LocationData struct {
	Latitude, Longitude float32 // required parameters
}

type OMessage struct {
	SendData                      // standard message options

	DisableWebPagePreview bool // pass true to suppress web page previews
}

type OPhoto struct {
	SendData
	MediaData
}

type OAudio struct {
	SendData
	MediaData
	LengthData
	ArtistData
}

type ODocument struct {
	SendData
	MediaData
}

type OSticker ODocument

type OVideo struct {
	SendData
	MediaData
	LengthData
	ResolutionData

	SupportsStreaming bool
}

type OAnimation struct {
	SendData
	MediaData
	LengthData
	ResolutionData
}

type OVoice struct {
	SendData
	MediaData
	LengthData
}

type OVideoNote struct {
	SendData
	MediaData
	LengthData
	ResolutionData
}

type OLocation struct {
	SendData
	LocationData

	LivePeriod int
}

type OVenue struct {
	SendData
	LocationData

	Title, Address string
	FoursquareId   string
	FoursquareType string
}

type OPoll struct {
	SendData

	Question              string
	Options             []string

	Type PollType
	IsAnonymous           bool
	AllowsMultipleAnswers bool
	IsClosed              bool

	// only required for quiz polls
	CorrectOptionId       int
}

type ODice struct {
	SendData
}

type OContact struct {
	SendData

	PhoneNumber string
	FirstName   string
	LastName    string
	VCard       string
}

type OChatAction struct {
	TargetData

	Action Status
}

type OChatMember struct {
	TargetData

	UserId UserID
}

type OMessageEdit struct {
	SourceData
	SendData

	DisableWebPagePreview bool
}

type OCaptionEdit struct {
	SourceData
	SendData
}

type OForward OCaptionEdit

type ODelete SourceData

type OStickerSet struct {
	Name string
}

type OInlineQueryAnswer struct {
	Id         InlineID
	Results  []interface{} // types: array of TInlineQueryResult*
	NextOffset string
	CacheTime  int
}

type OCallback struct {
	Id           CallbackID
	Notification string
	ShowAlert    bool
	CacheTime    int
	URL          string
}

type ORestrict struct {
	ChatID              interface{} // types: int, int64, or string
	UserID              int
	Until               int64
	CanSendMessages    *bool
	CanSendMedia       *bool
	CanSendPolls       *bool
	CanSendInline      *bool
	CanSendWebPreviews *bool
	CanChangeInfo      *bool
	CanInviteUsers     *bool
	CanPinMessages     *bool
}

type OGetFile struct {
	FileID string
}

type OFile struct {
	FilePath string
}
