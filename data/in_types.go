package data

import (
	"encoding/json"
)

type TUser struct {
	Id                int    `json:"id"`
	Bot               bool   `json:"is_bot"`
	First_name        string `json:"first_name"`
	Last_name        *string `json:"last_name"`
	Username         *string `json:"username"`
	Language_code    *string `json:"language_code"`
	Joins_groups     *bool   `json:"can_join_groups"`
	Privacy_disabled *bool   `json:"can_read_all_group_messages"`
	Inline_capable   *bool   `json:"supports_inline_queries"`
}

type TChat struct {
	Id                   int64            `json:"id"`
	Type                 string           `json:"type"`
	Title               *string           `json:"title"`
	Username            *string           `json:"username"`
	First_name          *string           `json:"first_name"`
	Last_name           *string           `json:"last_name"`
	Photo               *TChatPhoto       `json:"photo"`
	Description         *string           `json:"description"`
	Invite_link         *string           `json:"invite_link"`
	Pinned_message      *TMessage         `json:"pinned_message"`
	Permissions         *TChatPermissions `json:"permissions"`
	Slow_mode_delay     *int              `json:"slow_mode_delay"`
	Sticker_set_name    *string           `json:"sticker_set_name"`
	Can_set_sticker_set *bool             `json:"can_set_sticker_set"`
}

const Private string = "private"
const Group string = "group"
const Supergroup string = "supergroup"
const Channel string = "channel"

type TPhotoSize struct {
	File_id   string `json:"file_id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	File_size int    `json:"file_size"`
}

type TMaskPosition struct {
	Point	string  `json:"point"`
	X_shift float64 `json:"x_shift"`
	Y_shift float64 `json:"y_shift"`
	Scale   float64 `json:"scale"`
}

type TSticker struct {
	File_id        string        `json:"file_id"`
	Width          int           `json:"width"`
	Height         int           `json:"height"`
	Thumb         *TPhotoSize    `json:"thumb"`
	Emoji         *string        `json:"emoji"`
	Set_name      *string        `json:"set_name"`
	Mask_position *TMaskPosition `json:"mask_position"`
	File_size     *int           `json:"file_size"`
}

type TStickerSet struct {
	Name             string   `json:"name"`
	Title            string   `json:"title"`
	Contains_masks   bool     `json:"contains_masks"`
	Stickers      *[]TSticker `json:"stickers"`
}

type TMessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	User   TUser  `json:"user"`
}

type TMessage struct {
	Message_id          int              `json:"message_id"`
	From               *TUser            `json:"from"`
	Date                int64            `json:"date"`
	Chat                TChat            `json:"chat"`
	Forward_from       *TUser            `json:"forward_from"`
	Forward_from_chat  *TChat            `json:"forward_from_chat"`
	Forward_from_message_id *int         `json:"forward_from_message_id"`
	Forward_date       *int64            `json:"forward_date"`
	Reply_to_message   *TMessage         `json:"reply_to_message"`
	Edit_date          *int              `json:"edit_date"`
	Text               *string           `json:"text"`
	Caption            *string           `json:"caption"`
	Entities           *[]TMessageEntity `json:"entities"`
	Audio              *TGenericFile     `json:"audio"`
	Document           *TDocument        `json:"document"`
	Game               *TGame            `json:"game"`
	Photo              *[]TPhotoSize     `json:"photo"`
	Sticker            *TSticker         `json:"sticker"`
	Video              *TGenericFile     `json:"video"`
	Voice              *TGenericFile     `json:"voice"`
	Contact            *TContact         `json:"contact"`
	Location           *TLocation        `json:"location"`
	Venue              *TVenue           `json:"venue"`
	New_chat_member    *TUser            `json:"new_chat_member"`
	Left_chat_member   *TUser            `json:"left_chat_member"`
	New_chat_title     *string           `json:"new_chat_title"`
	New_chat_photo     *[]TPhotoSize     `json:"new_chat_photo"`
	Delete_chat_photo  *bool             `json:"delete_chat_photo"`
	Group_chat_created *bool             `json:"group_chat_created"`
	Supergroup_chat_created *bool        `json:"supergroup_chat_created"`
	Channel_chat_created    *bool        `json:"channel_chat_created"`
	Migrate_to_chat_id      *int64       `json:"migrate_to_chat_id"`
	Migrate_from_chat_id    *int64       `json:"migrate_from_chat_id"`
	Pinned_message     *TMessage         `json:"pinned_message"`
}

const Creator string = "creator"
const Admin string = "administrator"
const Member string = "member"
const Restricted string = "restricted"
const Left string = "left"
const Kicked string = "kicked"
const Banned string = Kicked

type TChatMember struct {
	User                   TUser  `json:"user"`
	Status                 string `json:"status"`
	// present only for restricted or kicked users
	Until_date            *int64  `json:"until_date"`
	// present only for administrators
	Can_be_edited         *bool   `json:"can_be_edited"`
	Can_change_info       *bool   `json:"can_change_info"`
	Can_post_messages     *bool   `json:"can_post_messages"`
	Can_edit_messages     *bool   `json:"can_edit_messages"`
	Can_delete_messages   *bool   `json:"can_delete_messages"`
	Can_invite_users      *bool   `json:"can_invite_users"`
	Can_restrict_members  *bool   `json:"can_restrict_members"`
	Can_pin_messages      *bool   `json:"can_pin_messages"`
	Can_promote_members   *bool   `json:"can_promote_members"`
	// present only for restricted users
	Is_member             *bool   `json:"is_member"`
	Can_send_anything     *bool   `json:"can_send_messages"`
	Can_send_media        *bool   `json:"can_send_media_messages"`
	Can_send_inline       *bool   `json:"can_send_other_messages"`
	Can_send_web_previews *bool   `json:"can_send_web_page_previews"`
}

type TGenericFile struct {
	File_id string `json:"file_id"`
}

type TFile struct {
	File_id    string `json:"file_id"`
	File_size *int    `json:"file_size"`
	File_path *string `json:"file_path"`
}

type TDocument struct {
	File_id    string     `json:"file_id"`
	Thumb     *TPhotoSize `json:"thumb"`
	File_name  string     `json:"file_name"`
	Mime_type  string     `json:"mime_type"`
	File_size  int        `json:"file_size"`
}

type TGame struct {
	Title string `json:"title"`
}

type TContact struct {
	Phone_number string `json:"phone_number"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	User_id      int    `json:"user_id"`
}

type TLocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type TVenue struct {
	Location    TLocation `json:"location"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type TInlineQuery struct {
	Id     string `json:"id"`
	From   TUser  `json:"from"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

type TChosenInlineResult struct {
	Result_id          string    `json:"result_id"`
	From               TUser     `json:"from"`
	Location          *TLocation `json:"location,omitempty"`
	Inline_message_id *string    `json:"inline_message_id,omitempty"`
	Query              string    `json:"query"`
}

type TCallbackQuery struct {
	Id                 string   `json:"id"`
	From               TUser    `json:"from"`
	Message           *TMessage `json:"message"`
	Inline_message_id *string   `json:"inline_message_id"`
	Chat_instance      string   `json:"chat_instance"`
	Data              *string   `json:"data"`
	Game_short_name   *string   `json:"game_short_name"`
}

type TUpdate struct {
	Update_id             int                 `json:"update_id"`
	Message              *TMessage            `json:"message,omitempty"`
	Edited_message       *TMessage            `json:"edited_message,omitempty"`
	Channel_post         *TMessage            `json:"channel_post,omitempty"`
	Edited_channel_post  *TMessage            `json:"edited_channel_post,omitempty"`
	Inline_query         *TInlineQuery        `json:"inline_query,omitempty"`
	Chosen_inline_result *TChosenInlineResult `json:"chosen_inline_result,omitempty"`
	Callback_query       *TCallbackQuery      `json:"callback_query,omitempty"`
}

type TGenericResponse struct {
	Ok          bool             `json:"ok"`
	Error_code  *int             `json:"error_code,omitempty"`
	Description *string          `json:"description,omitempty"`
	Result      *json.RawMessage `json:"result,omitempty"`
}

type TInlineQueryResultCachedSticker struct {
	Type                   string `json:"type"`
	Id                     string `json:"id"`
	Sticker_file_id        string `json:"sticker_file_id"`
	Reply_markup          *string `json:"reply_markup,omitempty"`
	Input_message_content *string `json:"input_message_content,omitempty"`
}

type TInlineQueryResultPhoto struct {
	Type                   string                   `json:"type"`
	Id                     string                   `json:"id"`
	Photo_url              string                   `json:"photo_url"`
	Thumb_url              string                   `json:"thumb_url"`
	Photo_width           *int                      `json:"photo_width,omitempty"`
	Photo_height          *int                      `json:"photo_height,omitempty"`
	Title                 *string                   `json:"title,omitempty"`
	Description           *string                   `json:"description,omitempty"`
	Caption               *string                   `json:"caption,omitempty"`
	Reply_markup          *string                   `json:"reply_markup,omitempty"`
	Input_message_content *TInputMessageTextContent `json:"input_message_content,omitempty"`
}

type TInlineQueryResultCachedPhoto struct {
	Type                   string                   `json:"type"`
	Id                     string                   `json:"id"`
	Photo_file_id          string                   `json:"photo_file_id"`
	Title                 *string                   `json:"title,omitempty"`
	Description           *string                   `json:"description,omitempty"`
	Caption               *string                   `json:"caption,omitempty"`
	Parse_mode            *string                   `json:"parse_mode,omitempty"`
	Reply_markup          *string                   `json:"reply_markup,omitempty"`
	Input_message_content *TInputMessageTextContent `json:"input_message_content,omitempty"`
}

type TInlineQueryResultGif struct {
	Type                   string                   `json:"type"`
	Id                     string                   `json:"id"`
	Gif_url                string                   `json:"gif_url"`
	Gif_width             *int                      `json:"gif_width,omitempty"`
	Gif_height            *int                      `json:"gif_height,omitempty"`
	Thumb_url              string                   `json:"thumb_url"`
	Title                 *string                   `json:"title,omitempty"`
	Caption               *string                   `json:"caption,omitempty"`
	Reply_markup          *string                   `json:"reply_markup,omitempty"`
	Input_message_content *TInputMessageTextContent `json:"input_message_content,omitempty"`
}

type TInlineKeyboard struct {
	Buttons [][]TInlineKeyboardButton `json:"inline_keyboard"`
}

type TInlineKeyboardButton struct {
	Text string `json:"text"`
	Data string `json:"callback_data"`
}

type TChatPermissions struct {
	Can_send_messages *bool `json:"can_send_messages,omitempty"`
	Can_send_media    *bool `json:"can_send_media_messages,omitempty"`
	Can_send_polls    *bool `json:"can_send_polls,omitempty"`
	Can_send_other    *bool `json:"can_send_other_messages,omitempty"`
	Can_preview_links *bool `json:"can_add_web_page_previews,omitempty"`
	Can_change_info   *bool `json:"can_change_info,omitempty"`
	Can_invite_users  *bool `json:"can_invite_users,omitempty"`
	Can_pin_messages  *bool `json:"can_pin_messages,omitempty"`
}

type TShippingAddress struct {
	Country_code string `json:"country_code"`
	State        string `json:"state"`
	City         string `json:"city"`
	Street_line1 string `json:"street_line1"`
	Street_line2 string `json:"street_line2"`
	Zip_code     string `json:"post_code"`
}

type TShippingQuery struct {
	Id               string           `json:"id"`
	From             TUser            `json:"from"`
	Invoice_payload  string           `json:"invoice_payload"`
	Shipping_address TShippingAddress `json:"shipping_address"`
}

type TOrderInfo struct {
	Name             *string           `json:"name,omitempty"`
	Phone_number     *string           `json:"phone_number,omitempty"`
	Email            *string           `json:"email,omitempty"`
	Shipping_address *TShippingAddress `json:"shipping_address,omitempty"`
}

type TPreCheckoutQuery struct {
	Id                  string     `json:"id"`
	From                TUser      `json:"from"`
	Currency            string     `json:"currency"`
	Total_amount        int64      `json:"total_amount"`
	Invoice_payload     string     `json:"invoice_payload"`
	Shipping_option_id *string     `json:"shipping_option_id"`
	Order_info         *TOrderInfo `json:"order_info"`
}

type TPollOption struct {
	Text  string `json:"text"`
	Votes int64  `json:"voter_count"`
}

type TPoll struct {
	Id             string      `json:"id"`
	Question       string      `json:"question"`
	Options      []TPollOption `json:"options"`
	Total_votes    int64       `json:"total_voter_count"`
	Closed         bool        `json:"is_closed"`
	Anonymous      bool        `json:"is_anonymous"`
	Type           string      `json:"type"`
	Multi_answer   bool        `json:"allows_multiple_answers"`
	Correct_answer int         `json:"correct_option_id"`
}

type TPollAnswer struct {
	Poll_id    string `json:"poll_id"`
	User       TUser  `json:"user"`
	Selected []int    `json:"option_ids"`
}

type TWebhookInfo struct {
	URL                 string `json:"url"`
	Custom_certificate  string `json:"has_custom_certificate"`
	Pending_updates     int    `json:"pending_update_count"`
	Last_error_date    *int    `json:"last_error_date,omitempty"`
	Last_error_message *string `json:"last_error_message,omitempty"`
	Max_connections    *int    `json:"max_connections"`
	Allowed_updates   []string `json:"allowed_updates"`
}

func (this *TInlineKeyboard) AddButton(b TInlineKeyboardButton) {
	if this.Buttons == nil { this.AddRow() }
	this.Buttons[len(this.Buttons) - 1] = append(this.Buttons[len(this.Buttons) - 1], b)
}

func (this *TInlineKeyboard) AddRow() {
	this.Buttons = append(this.Buttons, nil)
}

type TInputMessageTextContent struct {
	Message_text string `json:"message_text"`
	Parse_mode  *string `json:"parse_mode,omitempty"`
	No_preview  *bool   `json:"disable_web_page_preview,omitempty"`
}

type TChatPhoto struct {
	Small_id string `json:"small_file_id"`
	Large_id string `json:"big_file_id"`
}
