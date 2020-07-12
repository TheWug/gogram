package data

import (
	"strconv"
)

// unique identifiers used by telegram. treat these like atoms.
type UpdateID     int
type UserID       int
type ChatID       int64
type MsgID        int

func (this UpdateID)     String() string { return strconv.FormatInt(int64(this), 10) }
func (this UserID)       String() string { return strconv.FormatInt(int64(this), 10) }
func (this ChatID)       String() string { return strconv.FormatInt(int64(this), 10) }
func (this MsgID)        String() string { return strconv.FormatInt(int64(this), 10) }

type PollID       string
type FileID       string
type InlineID     string
type CallbackID   string
type ShippingID   string
type CheckoutID   string
type TxIDTelegram string
type TxIDVendor   string

func (this PollID)       String() string { return string(this) }
func (this FileID)       String() string { return string(this) }
func (this InlineID)     String() string { return string(this) }
func (this CallbackID)   String() string { return string(this) }
func (this ShippingID)   String() string { return string(this) }
func (this CheckoutID)   String() string { return string(this) }
func (this TxIDTelegram) String() string { return string(this) }
func (this TxIDVendor)   String() string { return string(this) }

// internal helper types
type DialogID     string

// enumerable values for TChat.Type
type ChatType string
const Private    ChatType = "private"
const Group      ChatType = "group"
const Supergroup ChatType = "supergroup"
const Channel    ChatType = "channel"
func (this ChatType) String() string { return string(this) }

// enumerable values for Outgoing.ParseMode
type MessageParseMode string
const ParseHTML      MessageParseMode = "HTML"
const ParseMarkdown  MessageParseMode = "Markdown"
const ParseMarkdown2 MessageParseMode = "Markdown2"
const ParseDefault   MessageParseMode = ""
func (this MessageParseMode) String() string { return string(this) }

// enumerable values for Outgoing.DisableNotification
type MessageNotify bool
const DisabledNo  MessageNotify = false
const DisabledYes MessageNotify = true
func (this MessageNotify) String() string { return strconv.FormatBool(bool(this)) }

// enumerable values for OChatAction.Action
type Status string
const Typing          Status = "typing"
const UploadPhoto     Status = "upload_photo"
const RecordVideo     Status = "record_video"
const UploadVideo     Status = "upload_video"
const RecordAudio     Status = "record_audio"
const UploadAudio     Status = "upload_audio"
const UploadDocument  Status = "upload_document"
const FindLocation    Status = "find_location"
const RecordVideoNote Status = "record_video_note"
const UploadVideoNote Status = "upload_video_note"
func (this Status) String() string { return string(this) }

// enumerable values for TMessageEntity.Type
type EntityType string
const Mention     EntityType = "mention"
const TextMention EntityType = "text_mention"
const Hashtag     EntityType = "hashtag"
const Cashtag     EntityType = "cashtag"
const Command     EntityType = "bot_command"
const URL         EntityType = "url"
const Email       EntityType = "email"
const Phone       EntityType = "phone_number"
const Bold        EntityType = "bold"
const Italic      EntityType = "italic"
const Underline   EntityType = "underline"
const Strike      EntityType = "strikethrough"
const Code        EntityType = "code"
const Pre         EntityType = "pre"
const TextLink    EntityType = "text_link"
func (this EntityType) String() string { return string(this) }

// enumerable values for TChatMember.Status
type MemberStatus string
const Creator    MemberStatus = "creator"
const Admin      MemberStatus = "administrator"
const Member     MemberStatus = "member"
const Restricted MemberStatus = "restricted"
const Left       MemberStatus = "left"
const Kicked     MemberStatus = "kicked"
const Banned     MemberStatus = Kicked
func (this MemberStatus) String() string { return string(this) }

// enumerable values for OPoll.Type, TPoll.Type
type PollType string
const NormalPoll PollType = "regular"
const QuizPoll   PollType = "quiz"
func (this PollType) String() string { return string(this) }

// enumerable values for TMaskPosition.Point
type MaskPoint string
const Forehead MaskPoint = "forehead"
const Eyes     MaskPoint = "eyes"
const Mouth    MaskPoint = "mouth"
const Chin     MaskPoint = "chin"
func (this MaskPoint) String() string { return string(this) }
