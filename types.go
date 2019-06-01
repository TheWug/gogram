package telegram

import (
	"encoding/json"
)

type UserID int
type ChatID int64

type Sender struct {
	User UserID
	Channel ChatID
}

type TUser struct {
	Id          int    `json:"id"`
	First_name  string `json:"first_name"`
	Last_name  *string `json:"last_name"`
	Username   *string `json:"username"`
}

type TChat struct {
	Id          int64  `json:"id"`
	Type        string `json:"type"`
	Title      *string `json:"title"`
	Username   *string `json:"username"`
	First_name *string `json:"first_name"`
	Last_name  *string `json:"last_name"`
}

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

type TChatMember struct {
	User                  TUser  `json:"user"`
	Status                string `json:"status"`
	Can_restrict_members *bool   `json:"can_restrict_members"`
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

func (this *TInlineKeyboard) AddButton(b TInlineKeyboardButton) {
	if this.Buttons == nil { this.AddRow() }
	this.Buttons[len(this.Buttons) - 1] = append(this.Buttons[len(this.Buttons) - 1], b)
}

func (this *TInlineKeyboard) AddRow() {
	this.Buttons = append(this.Buttons, nil)
}

type TInputMessageTextContent struct {
	Message_text string `json:"message_text"`
	Parse_mode   string `json:"parse_mode"`
	No_preview   bool   `json:"disable_web_page_preview"`
}

type ResponseHandler interface {
	Callback(result *json.RawMessage, success bool, err error, http_code int)
}
