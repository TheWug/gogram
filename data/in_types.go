package data

import (
	"encoding/json"
)

type TUser struct {
	Id                       UserID `json:"id"`
	IsBot                    bool   `json:"is_bot"`
	FirstName                string `json:"first_name"`
	LastName                *string `json:"last_name"`
	Username                *string `json:"username"`
	LanguageCode            *string `json:"language_code"`
	CanJoinGroups           *bool   `json:"can_join_groups"`
	CanReadAllGroupMessages *bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   *bool   `json:"supports_inline_queries"`
}

type TChat struct {
	Id                ChatID           `json:"id"`
	Type              ChatType         `json:"type"`
	Title            *string           `json:"title"`
	Username         *string           `json:"username"`
	FirstName        *string           `json:"first_name"`
	LastName         *string           `json:"last_name"`
	Photo            *TChatPhoto       `json:"photo"`
	Description      *string           `json:"description"`
	InviteLink       *string           `json:"invite_link"`
	PinnedMessage    *TMessage         `json:"pinned_message"`
	Permissions      *TChatPermissions `json:"permissions"`
	SlowModeDelay    *int              `json:"slow_mode_delay"`
	StickerSetName   *string           `json:"sticker_set_name"`
	CanSetStickerSet *bool             `json:"can_set_sticker_set"`
}

type TPhotoSize struct {
	Id       FileID `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type TMaskPosition struct {
	Point	MaskPoint `json:"point"`
	X_Shift float64   `json:"x_shift"`
	Y_Shift float64   `json:"y_shift"`
	Scale   float64   `json:"scale"`
}

type TSticker struct {
	Id            FileID        `json:"file_id"`
	Width         int           `json:"width"`
	Height        int           `json:"height"`
	Animated      bool          `json:"is_animated"`
	Thumb        *TPhotoSize    `json:"thumb"`
	Emoji        *string        `json:"emoji"`
	SetName      *string        `json:"set_name"`
	MaskPosition *TMaskPosition `json:"mask_position"`
	FileSize     *int           `json:"file_size"`
}

type TAudio struct {
	Id         FileID     `json:"file_id"`
	UniqueId   FileID     `json:"unique_file_id"`
	Duration   int        `json:"duration"`
	Performer *string     `json:"performer"`
	Title     *string     `json:"title"`
	MimeType  *string     `json:"mime_type"`
	FileSize  *int        `json:"file_size"`
	Thumb     *TPhotoSize `json:"thumb"`
}

type TVoice struct {
	Id        FileID `json:"file_id"`
	UniqueId  string `json:"unique_file_id"`
	Duration  int    `json:"duration"`
	MimeType *string `json:"mime_type"`
	FileSize *int    `json:"file_size"`
}

type TVideo struct {
	Id        FileID     `json:"file_id"`
	UniqueId  FileID     `json:"unique_file_id"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Duration  int        `json:"duration"`
	Thumb    *TPhotoSize `json:"thumb"`
	MimeType *string     `json:"mime_type"`
	FileSize *int        `json:"file_size"`
}

type TVideoNote struct {
	Id        FileID `json:"file_id"`
	UniqueId  FileID `json:"unique_file_id"`
	Diameter  int    `json:"length"`
	Duration  int    `json:"duration"`
	MimeType *string `json:"mime_type"`
	FileSize *int    `json:"file_size"`
}

type TStickerSet struct {
	Name          string     `json:"name"`
	Title         string     `json:"title"`
	Animated      bool       `json:"is_animated"`
	ContainsMasks bool       `json:"contains_masks"`
	Stickers   *[]TSticker   `json:"stickers"`
	Thumb         TPhotoSize `json:"thumb"`
}

type TMessageEntity struct {
	Type      EntityType `json:"type"`
	Offset    int        `json:"offset"`
	Length    int        `json:"length"`
	Url      *string     `json:"url"`
	User     *TUser      `json:"user"`
	Language *string     `json:"language"`
}

type TMessage struct {
	Id                     MsgID              `json:"message_id"`
	From                  *TUser              `json:"from"`
	Date                   int64              `json:"date"`
	Chat                   TChat              `json:"chat"`
	ForwardFrom           *TUser              `json:"forward_from"`
	ForwardFromChat       *TChat              `json:"forward_from_chat"`
	ForwardFromMessageId  *int                `json:"forward_from_message_id"`
	ForwardSignature      *string             `json:"forward_signature"`
	ForwardSender         *string             `json:"forward_sender_name"`
	ForwardDate           *int64              `json:"forward_date"`
	ReplyToMessage        *TMessage           `json:"reply_to_message"`
	EditDate              *int                `json:"edit_date"`
	MediaGroupId          *string             `json:"media_group_id"`
	AuthorSignature       *string             `json:"author_signature"`
	Text                  *string             `json:"text"`
	Entities            *[]TMessageEntity     `json:"entities"`
	Caption               *string             `json:"caption"`
	CaptionEntities     *[]TMessageEntity     `json:"caption_entities"`
	Audio                 *TAudio             `json:"audio"`
	Document              *TDocument          `json:"document"`
	Animation             *TAnimation         `json:"animation"`
	Game                  *TGame              `json:"game"`
	Photo                 *[]TPhotoSize       `json:"photo"`
	Sticker               *TSticker           `json:"sticker"`
	Video                 *TVideo             `json:"video"`
	Voice                 *TVoice             `json:"voice"`
	Video_note            *TVideoNote         `json:"video_note"`
	Contact               *TContact           `json:"contact"`
	Location              *TLocation          `json:"location"`
	Venue                 *TVenue             `json:"venue"`
	Poll                  *TPoll              `json:"poll"`
	Dice                  *TDice              `json:"dice"`
	NewChatMembers      *[]TUser              `json:"new_chat_member"`
	LeftChatMember        *TUser              `json:"left_chat_member"`
	NewChatTitle          *string             `json:"new_chat_title"`
	NewChatPhoto        *[]TPhotoSize         `json:"new_chat_photo"`
	DeleteChatPhoto       *bool               `json:"delete_chat_photo"`
	GroupChatCreated      *bool               `json:"group_chat_created"`
	SupergroupChatCreated *bool               `json:"supergroup_chat_created"`
	ChannelChatCreated    *bool               `json:"channel_chat_created"`
	MigrateToChatId       *int64              `json:"migrate_to_chat_id"`
	MigrateFromChatId     *int64              `json:"migrate_from_chat_id"`
	PinnedMessage         *TMessage           `json:"pinned_message"`
	Invoice               *TInvoice           `json:"invoice"`
	SuccessfulPayment     *TSuccessfulPayment `json:"successful_payment"`
	ConnectedWebsite      *string             `json:"connected_website"`
	PassportData          *TPassportData      `json:"passport_data"`
	ReplyMarkup           *TInlineKeyboard    `json:"reply_markup"`
}

type TChatMember struct {
	User                TUser        `json:"user"`
	Status              MemberStatus `json:"status"`
	// present only for restricted or kicked users
	UntilDate          *int64        `json:"until_date"`
	// present only for administrators
	CanBeEdited        *bool         `json:"can_be_edited"`
	CanChangeInfo      *bool         `json:"can_change_info"`
	CanPostMessages    *bool         `json:"can_post_messages"`
	CanEditMessages    *bool         `json:"can_edit_messages"`
	CanDeleteMessages  *bool         `json:"can_delete_messages"`
	CanInviteUsers     *bool         `json:"can_invite_users"`
	CanRestrictMembers *bool         `json:"can_restrict_members"`
	CanPinMessages     *bool         `json:"can_pin_messages"`
	CanPromoteMembers  *bool         `json:"can_promote_members"`
	// present only for restricted users
	IsMember           *bool         `json:"is_member"`
	CanSendMessages    *bool         `json:"can_send_messages"`          // can they send anything at all
	CanSendMedia       *bool         `json:"can_send_media_messages"`    // can they send uploadable media
	CanSendInline      *bool         `json:"can_send_other_messages"`    // can they send stickers, gifs, or use inline bots
	CanSendWebPreviews *bool         `json:"can_send_web_page_previews"` // can they attach webpage previews to their messages
}

type TGenericFile struct {
	Id       FileID `json:"file_id"`
	UniqueId FileID `json:"file_unique_id"` // can be used to compare files between bots, but can't be used to send or download messages
}

type TFile struct {
	Id        FileID `json:"file_id"`
	FileSize *int    `json:"file_size"`
	FilePath *string `json:"file_path"`
}

type TDocument struct {
	Id        FileID     `json:"file_id"`
	Thumb    *TPhotoSize `json:"thumb"`
	FileName *string     `json:"file_name"`
	MimeType *string     `json:"mime_type"`
	FileSize *int        `json:"file_size"`
}

type TAnimation struct {
	Id        FileID     `json:"file_id"`
	UniqueId  FileID     `json:"file_unique_id"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Duration  int        `json:"duration"`
	Thumb    *TPhotoSize `json:"thumb"`
	FileName *string     `json:"file_name"`
	MimeType *string     `json:"mime_type"`
	FileSize *int        `json:"file_size"`
}

type TGame struct {
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Photo         []TPhotoSize     `json:"photo"`
	Text           *string         `json:"text"`
	TextEntities *[]TMessageEntity `json:"text_entities"`
	Animation      *TAnimation     `json:"animation"`
}

type TGameHighScore struct {
	Position int   `json:"position"`
	User     TUser `json:"user"`
	Score    int64 `json:"score"`
}

type TContact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName   *string `json:"last_name,omitempty"`
	UserID     *int    `json:"user_id,omitempty"`
	Vcard      *string `json:"vcard,omitempty"`
}

type TLocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type TVenue struct {
	Location        TLocation `json:"location"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	FoursquareId   *string    `json:"foursquare_id"`
	FoursquareType *string    `json:"foursquare_type"`
}

type TDice struct {
	Value int `json:"value"`
}

type TInlineQuery struct {
	Id     InlineID `json:"id"`
	From   TUser    `json:"from"`
	Query  string   `json:"query"`
	Offset string   `json:"offset"`
}

type TChosenInlineResult struct {
	ResultId         string    `json:"result_id"`
	From             TUser     `json:"from"`
	Location        *TLocation `json:"location,omitempty"`
	InlineMessageId *InlineID  `json:"inline_message_id,omitempty"`
	Query            string    `json:"query"`
}

type TCallbackQuery struct {
	Id               CallbackID `json:"id"`
	From             TUser      `json:"from"`
	Message         *TMessage   `json:"message"`
	InlineMessageId *string     `json:"inline_message_id"`
	ChatInstance     string     `json:"chat_instance"`
	Data            *string     `json:"data"`
	GameShortName   *string     `json:"game_short_name"`
}

type TUpdate struct {
	Id                  UpdateID            `json:"update_id"`
	Message            *TMessage            `json:"message,omitempty"`
	EditedMessage      *TMessage            `json:"edited_message,omitempty"`
	ChannelPost        *TMessage            `json:"channel_post,omitempty"`
	EditedChannelPost  *TMessage            `json:"edited_channel_post,omitempty"`
	InlineQuery        *TInlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *TChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery      *TCallbackQuery      `json:"callback_query,omitempty"`
}

type TGenericResponse struct {
	Ok           bool            `json:"ok"`
	ErrorCode   *int             `json:"error_code,omitempty"`
	Description *string          `json:"description,omitempty"`
	Result      *json.RawMessage `json:"result,omitempty"`
}

type TInlineQueryResultCachedSticker struct {
	Type                 string `json:"type"`
	Id                   string `json:"id"`
	StickerId            FileID `json:"sticker_file_id"`
	ReplyMarkup         *string `json:"reply_markup,omitempty"`
	InputMessageContent *string `json:"input_message_content,omitempty"`
}

type TInlineQueryResultPhoto struct {
	Type                 string                   `json:"type"`
	Id                   string                   `json:"id"`
	PhotoUrl             string                   `json:"photo_url"`
	ThumbUrl             string                   `json:"thumb_url"`
	PhotoWidth          *int                      `json:"photo_width,omitempty"`
	PhotoHeight         *int                      `json:"photo_height,omitempty"`
	Title               *string                   `json:"title,omitempty"`
	Description         *string                   `json:"description,omitempty"`
	Caption             *string                   `json:"caption,omitempty"`
	ReplyMarkup         *string                   `json:"reply_markup,omitempty"`
	InputMessageContent *TInputMessageTextContent `json:"input_message_content,omitempty"`
}

type TInlineQueryResultCachedPhoto struct {
	Type                 string                   `json:"type"`
	Id                   string                   `json:"id"`
	PhotoId              FileID                   `json:"photo_file_id"`
	Title               *string                   `json:"title,omitempty"`
	Description         *string                   `json:"description,omitempty"`
	Caption             *string                   `json:"caption,omitempty"`
	ParseMode           *string                   `json:"parse_mode,omitempty"`
	ReplyMarkup         *string                   `json:"reply_markup,omitempty"`
	InputMessageContent *TInputMessageTextContent `json:"input_message_content,omitempty"`
}

type TInlineQueryResultGif struct {
	Type                 string                   `json:"type"`
	Id                   string                   `json:"id"`
	GifUrl               string                   `json:"gif_url"`
	GifWidth            *int                      `json:"gif_width,omitempty"`
	GifHeight           *int                      `json:"gif_height,omitempty"`
	ThumbUrl            string                   `json:"thumb_url"`
	Title               *string                   `json:"title,omitempty"`
	Caption             *string                   `json:"caption,omitempty"`
	ParseMode           *string                   `json:"parse_mode,omitempty"`
	ReplyMarkup         *string                   `json:"reply_markup,omitempty"`
	InputMessageContent *TInputMessageTextContent `json:"input_message_content,omitempty"`
}

type TInlineKeyboard struct {
	Buttons [][]TInlineKeyboardButton `json:"inline_keyboard"`
}

type TInlineKeyboardButton struct {
	Text              string        `json:"text"`
	Url              *string        `json:"url,omitempty"`
	LoginUrl         *TLoginURL     `json:"login_url,omitempty"`
	Data             *string        `json:"callback_data,omitempty"`
	SwitchInline     *string        `json:"switch_inline_query,omitempty"`
	SwitchInlineHere *string        `json:"switch_inline_query_current_chat,omitempty"`
	CallbackGame     *TCallbackGame `json:"callback_game,omitempty"`
	Pay              *bool          `json:"pay,omitempty"`
}

type TReplyKeyboard struct {
	Buttons [][]TKeyboardButton `json:"keyboard"`
	Resizable  *bool            `json:"resize_keyboard,omitempty"`
	OneTime    *bool            `json:"one_time_keyboard,omitempty"`
	Selective  *bool            `json:"selective,omitempty"`
}

type TKeyboardButton struct {
	Text             string                  `json:"text"`
	RequestContact  *bool                    `json:"request_contact,omitempty"`
	RequestLocation *bool                    `json:"request_location,omitempty"`
	RequestPoll     *TKeyboardButtonPollType `json:"request_poll,omitempty"`
}

type TKeyboardButtonPollType struct {
	Type PollType `json:"type"`
}

type TReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective     *bool `json:"selective,omitempty"`
}

type TForceReply struct {
	ForceReply bool `json:"force_reply"`
	Selective *bool `json:"selective,omitempty"`
}

type TLoginURL struct {
	Url            string `json:"url"`
	ForwardText   *string `json:"forward_text"`
	BotUsername   *string `json:"bot_username"`
	MsgPermission *bool   `json:"request_write_access"`
}

type TCallbackGame struct {
	// contains no fields.
}

type TChatPermissions struct {
	CanSendMessages *bool `json:"can_send_messages,omitempty"`
	CanSendMedia    *bool `json:"can_send_media_messages,omitempty"`
	CanSendPolls    *bool `json:"can_send_polls,omitempty"`
	CanSendOther    *bool `json:"can_send_other_messages,omitempty"`
	CanPreviewLinks *bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo   *bool `json:"can_change_info,omitempty"`
	CanInviteUsers  *bool `json:"can_invite_users,omitempty"`
	CanPinMessages  *bool `json:"can_pin_messages,omitempty"`
}

type TShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	ZipCode     string `json:"post_code"`
}

type TShippingQuery struct {
	Id              ShippingID       `json:"id"`
	From            TUser            `json:"from"`
	InvoicePayload  string           `json:"invoice_payload"`
	ShippingAddress TShippingAddress `json:"shipping_address"`
}

type TOrderInfo struct {
	Name            *string           `json:"name,omitempty"`
	PhoneNumber     *string           `json:"phone_number,omitempty"`
	Email           *string           `json:"email,omitempty"`
	ShippingAddress *TShippingAddress `json:"shipping_address,omitempty"`
}

type TInvoice struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Cookie      string `json:"start_parameter"`
	Currency    string `json:"currency"`
	Total       int64  `json:"total_amount"`
}

type TPreCheckoutQuery struct {
	Id                CheckoutID `json:"id"`
	From              TUser      `json:"from"`
	Currency          string     `json:"currency"`
	TotalAmount       int64      `json:"total_amount"`
	InvoicePayload    string     `json:"invoice_payload"`
	ShippingOptionId *string     `json:"shipping_option_id"`
	OrderInfo        *TOrderInfo `json:"order_info"`
}

type TSuccessfulPayment struct {
	Currency          string       `json:"currency"`
	Total             int64        `json:"total_amount"`
	InvoicePayload    string       `json:"invoice_payload"`
	ShippingOptionId *string       `json:"shipping_option_id"`
	OrderInfo        *TOrderInfo   `json:"order_info"`
	TxIdTelegram      TxIDTelegram `json:"telegram_payment_charge_id"`
	TxIdProvider      TxIDVendor   `json:"provider_payment_charge_id"`
}

type TPollOption struct {
	Text  string `json:"text"`
	VoterCount int64  `json:"voter_count"`
}

type TPoll struct {
	Id                    PollID      `json:"id"`
	Question              string      `json:"question"`
	Options             []TPollOption `json:"options"`
	TotalVoterCount       int64       `json:"total_voter_count"`
	IsClosed              bool        `json:"is_closed"`
	IsAnonymous           bool        `json:"is_anonymous"`
	Type                  string      `json:"type"`
	AllowsMultipleAnswers bool        `json:"allows_multiple_answers"`
	CorrectOptionId      *int         `json:"correct_option_id"`
}

type TPollAnswer struct {
	Id         PollID `json:"poll_id"`
	User       TUser  `json:"user"`
	Selected []int    `json:"option_ids"`
}

type TPassportData struct {
	// stub. implement more of this later.
}

type TWebhookInfo struct {
	URL                  string `json:"url"`
	HasCustomCertificate string `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
	LastErrorDate       *int    `json:"last_error_date,omitempty"`
	LastErrorMessage    *string `json:"last_error_message,omitempty"`
	MaxConnections      *int    `json:"max_connections"`
	AllowedUpdates     []string `json:"allowed_updates"`
}

func (this *TInlineKeyboard) AddButton(b TInlineKeyboardButton) {
	if this.Buttons == nil { this.AddRow() }
	this.Buttons[len(this.Buttons) - 1] = append(this.Buttons[len(this.Buttons) - 1], b)
}

func (this *TInlineKeyboard) AddRow() {
	this.Buttons = append(this.Buttons, nil)
}

type TInputMessageTextContent struct {
	MessageText string `json:"message_text"`
	ParseMode  *string `json:"parse_mode,omitempty"`
	NoPreview  *bool   `json:"disable_web_page_preview,omitempty"`
}

type TChatPhoto struct {
	SmallId FileID `json:"small_file_id"`
	LargeId FileID `json:"big_file_id"`
}
