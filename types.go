package telegram

import (
	"encoding/json"
)

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

type TSticker struct {
	File_id    string     `json:"file_id"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
//	Thumb     *TPhotosize `json:"thumb"`
	Emoji     *string     `json:"emoji"`
	File_size *int        `json:"file_size"`
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
	Forward_date       *int64            `json:"forward_date"`
	Reply_to_message   *TMessage         `json:"reply_to_message"`
	Edit_date          *int              `json:"edit_date"`
	Text               *string           `json:"text"`
	Caption            *string           `json:"caption"`
	Entities           *[]TMessageEntity `json:"entities"`
	Audio              *TGenericFile     `json:"audio"`
	Document           *TGenericFile     `json:"document"`
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

type TGenericFile struct {
	File_id string `json:"file_id"`
}

type TGame struct {
	Title string `json:"title"`
}

type TPhotoSize struct {
	File_id   string `json:"file_id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	File_size int    `json:"file_size"`
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

type TUpdate struct {
	Update_id             int                 `json:"update_id"`
	Message              *TMessage            `json:"message,omitempty"`
	Edited_message       *TMessage            `json:"edited_message,omitempty"`
	Inline_query         *TInlineQuery        `json:"inline_query,omitempty"`
	Chosen_inline_result *TChosenInlineResult `json:"chosen_inline_result,omitempty"`
//	Callback_query       *TCallbackQuery      `json:"callback_query,omitempty"`
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
	Type                   string `json:"type"`
	Id                     string `json:"id"`
	Photo_url              string `json:"photo_url"`
	Thumb_url              string `json:"thumb_url"`
	Photo_width           *int    `json:"photo_width,omitempty"`
	Photo_height          *int    `json:"photo_height,omitempty"`
	Title                 *string `json:"title,omitempty"`
	Description           *string `json:"description,omitempty"`
	Caption               *string `json:"caption,omitempty"`
	Reply_markup          *string `json:"reply_markup,omitempty"`
	Input_message_content *string `json:"input_message_content,omitempty"`
}

type TInlineQueryResultGif struct {
	Type                   string `json:"type"`
	Id                     string `json:"id"`
	Gif_url                string `json:"gif_url"`
	Gif_width             *int    `json:"gif_width,omitempty"`
	Gif_height            *int    `json:"gif_height,omitempty"`
	Thumb_url              string `json:"thumb_url"`
	Title                 *string `json:"title,omitempty"`
	Caption               *string `json:"caption,omitempty"`
	Reply_markup          *string `json:"reply_markup,omitempty"`
	Input_message_content *string `json:"input_message_content,omitempty"`
}
