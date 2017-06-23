package telegram

import (
	"fmt"
	"strings"
	"time"
)

// width 10
func (this TUser) IdString() (string) {
	return fmt.Sprintf("U%-9d", this.Id)
}

// width 56
func (this TUser) PrintableString() (string) {
	var username  string = "<no username>"
	var name      string = ""

	if this.Username != nil { username = "@" + *this.Username }
	if this.Last_name != nil { name = this.First_name + " " + *this.Last_name }
	if this.Last_name == nil { name = this.First_name }
	name = strings.Trim(strings.Replace(name, " ", "-", -1), " \r\n\t")

	return fmt.Sprintf("%-56s", fmt.Sprintf("%s %s %s", this.IdString(), username, name))
}

func (this TUser) NameString() (string) {
	if this.Last_name != nil { return this.First_name + " " + *this.Last_name }
	return this.First_name
}

func (this TUser) UsernameString() (string) {
	if this.Username != nil { return *this.Username }
	return "(none)"
}

// width 16
func (this TChat) IdString() (string) {
	return fmt.Sprintf("C%-15d", this.Id)
}

// width variable
func (this TChat) TitleString() (string) {
	if this.Title == nil { return this.IdString() }
	return *this.Title
}

// width 64
func (this TChat) PrintableString() (string) {
	var title string = ""
	var username string = "<no username>"
	var msgtype = ' '

	if this.Title != nil { title = *this.Title }
	if this.First_name != nil { title = *this.First_name }
	if this.Last_name != nil { title = title + " " + *this.Last_name }
	title = strings.Trim(strings.Replace(title, " ", "-", -1), " \r\n\t")

	if this.Username != nil { username = "@" + *this.Username }

	if this.Type == "private" { msgtype = 'P' }
	if this.Type == "supergroup" { msgtype = 'S' }
	if this.Type == "group" { msgtype = 'G' }
	if this.Type == "channel" { msgtype = 'C' }

	return fmt.Sprintf("%-64s", fmt.Sprintf("%c%s %s %s", msgtype, this.IdString(), username, title))
}

// width 10
func (this *TMessage) IdString() (string) {
	return fmt.Sprintf("M%-9d", this.Message_id)
}

// width 11
func (this *TMessage) IdStringWithEdit(is_edit bool) (string) {
	edited := ' '

	if is_edit { edited = 'E' }

	return fmt.Sprintf("%cM%-9d", edited, this.Message_id)
}

// width 79
func (this *TMessage) RplFwdString() (string) {
	var intermediate string = "NML -          -          - -"
	if this.Forward_from != nil { intermediate = fmt.Sprintf("FWD -          %s", this.Forward_from.PrintableString()) }
	if this.Forward_from_chat != nil { intermediate = fmt.Sprintf("FWD -          %s", this.Forward_from_chat.PrintableString()) }

	if this.Reply_to_message != nil && this.Reply_to_message.From != nil { intermediate = fmt.Sprintf("RPL %s %s", this.Reply_to_message.IdString(), this.Reply_to_message.From.PrintableString()) }
	if this.Reply_to_message != nil && this.Reply_to_message.From == nil { intermediate = fmt.Sprintf("RPL %s %s", this.Reply_to_message.IdString(), this.Reply_to_message.Chat.PrintableString()) } // longest, width 79

	return fmt.Sprintf("%-79s", intermediate)
}

// width: 8
func (this *TMessage) TypeString() (string) {
	var intermediate string = "UNKNOWN"
	if this.Text != nil { intermediate = "MESSAGE"
	} else if this.Audio != nil { intermediate = "AUDIO"
	} else if this.Document != nil { intermediate = "DOCUMENT"
	} else if this.Game != nil { intermediate = "GAME"
	} else if this.Photo != nil { intermediate = "PHOTO"
	} else if this.Sticker != nil { intermediate = "STICKER"
	} else if this.Video != nil { intermediate = "VIDEO"
	} else if this.Voice != nil { intermediate = "VOICE"
	} else if this.Contact != nil { intermediate = "CONTACT"
	} else if this.Location != nil { intermediate = "LOCATION"
	} else if this.Venue != nil { intermediate = "VENUE"
	} else if this.New_chat_member != nil { intermediate = "ADDUSER"
	} else if this.Left_chat_member != nil { intermediate = "DELUSER"
	} else if this.New_chat_title != nil { intermediate = "SETTITLE"
	} else if this.New_chat_photo != nil { intermediate = "SETPHOTO"
	} else if this.Delete_chat_photo != nil { intermediate = "DELPHOTO"
	} else if this.Group_chat_created != nil { intermediate = "MKGROUP"
	} else if this.Supergroup_chat_created != nil { intermediate = "MKSGROUP"
	} else if this.Channel_chat_created != nil { intermediate = "MKCHANNL"
	} else if this.Migrate_to_chat_id != nil { intermediate = "TOSGRP"
	} else if this.Migrate_from_chat_id != nil { intermediate = "FROMSGRP"
	} else if this.Pinned_message != nil { intermediate = "PINMSG" }
	return fmt.Sprintf("%-8s", intermediate)
}

// width: 41
func (this *TMessage) TimestampString() (string) {
	tf_string := "2006-01-02 15:04:05"
	if this.Forward_date != nil {
		return fmt.Sprintf("%s (%s)", time.Unix(this.Date, 0).UTC().Format(tf_string), time.Unix(*this.Forward_date, 0).UTC().Format(tf_string))
	} else {
		return fmt.Sprintf("%-41s", time.Unix(this.Date, 0).UTC().Format(tf_string))
	}
}

// wide as fuck
func (this *TMessage) PrintableString(is_edit bool) (string) {
	var message_ts, message_type, message_info, fwdrpl_info, sender_info, receiver_info, message_contents string

	message_ts = this.TimestampString()
	message_type = this.TypeString()
	message_info = this.IdStringWithEdit(is_edit)
	fwdrpl_info = this.RplFwdString()

	if this.From != nil { sender_info = this.From.PrintableString() }
	if this.From == nil { sender_info = this.Chat.PrintableString() }

	receiver_info = this.Chat.PrintableString()

	message_contents = this.MessageContents()
	
	return fmt.Sprintf("%s %s %s %s %s %s %s", message_ts, message_type, message_info, sender_info, receiver_info, fwdrpl_info, message_contents)
}

// width: variable.
func (this *TMessage) MessageContents() (string) {
	var caption string = "(no caption)"
	if this.Caption != nil { caption = fmt.Sprintf("<%s>", strings.Replace(*this.Caption, "\n", "", -1)) }

	if this.Audio != nil { return this.Audio.File_id }
	if this.Document != nil { return fmt.Sprintf("%s %s", this.Document.File_id, caption) }
	if this.Game != nil { return this.Game.Title }
	if this.Photo != nil { return fmt.Sprintf("%s %s", GetLargestPhoto(this.Photo).File_id, caption) }
	if this.Sticker != nil { return this.Sticker.File_id }
	if this.Video != nil { return fmt.Sprintf("%s %s", this.Video.File_id, caption) }
	if this.Voice != nil { return this.Voice.File_id }
	if this.Contact != nil { return this.Contact.Phone_number + " " + this.Contact.First_name }
	if this.Location != nil { return fmt.Sprintf("(%f %f)", this.Location.Longitude, this.Location.Latitude) }
	if this.Venue != nil { return fmt.Sprintf("(%f %f) %s: %s", this.Venue.Location.Longitude, this.Venue.Location.Latitude, this.Venue.Title, this.Venue.Description) }
	if this.New_chat_member != nil { return fmt.Sprintf("[Chat member added] %s", this.New_chat_member.PrintableString()) }
	if this.Left_chat_member != nil { return fmt.Sprintf("[Chat member removed] %s", this.Left_chat_member.PrintableString()) }
	if this.New_chat_title != nil { return fmt.Sprintf("[New chat title] %s", *this.New_chat_title) }
	if this.New_chat_photo != nil { return fmt.Sprintf("[Chat photo set] %s", GetLargestPhoto(this.New_chat_photo).File_id) }
	if this.Delete_chat_photo != nil { return "[Chat photo deleted]" }
	if this.Group_chat_created != nil { return "[Group created]" }
	if this.Supergroup_chat_created != nil { return "[Supergroup created]" }
	if this.Channel_chat_created != nil { return "[Channel created]" }
	if this.Migrate_to_chat_id != nil { return fmt.Sprintf("[supergroup created from existing chat] %d", this.Migrate_to_chat_id) }
	if this.Migrate_from_chat_id != nil { return fmt.Sprintf("[Group converted to supergroup] %d", this.Migrate_from_chat_id) }
	if this.Pinned_message != nil { return fmt.Sprintf("[Message pinned] %s", this.Pinned_message.MessageContents()) }
	if this.Text != nil { return strings.Replace(*this.Text, "\n", "", -1) }
	return "Unknown"
}

func GetLargestPhoto(this *[]TPhotoSize) (*TPhotoSize) {
	var largest *TPhotoSize = nil
	var largest_size int64 = 0

	if this == nil { return nil }

	for _, item := range *this {
		if int64(item.Width) * int64(item.Height) > largest_size { largest = &item }
	}

	return largest
}
