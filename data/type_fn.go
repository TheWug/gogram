package data

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
	if this.LastName != nil { name = this.FirstName + " " + *this.LastName }
	if this.LastName == nil { name = this.FirstName }
	name = strings.Trim(strings.Replace(name, " ", "-", -1), " \r\n\t")

	return fmt.Sprintf("%-56s", fmt.Sprintf("%s %s %s", this.IdString(), username, name))
}

func (this TUser) NameString() (string) {
	if this.LastName != nil { return this.FirstName + " " + *this.LastName }
	return this.FirstName
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
	if this.FirstName != nil { title = *this.FirstName }
	if this.LastName != nil { title = title + " " + *this.LastName }
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
	return fmt.Sprintf("M%-9d", this.Id)
}

// width 11
func (this *TMessage) IdStringWithEdit(is_edit bool) (string) {
	edited := ' '

	if is_edit { edited = 'E' }

	return fmt.Sprintf("%cM%-9d", edited, this.Id)
}

// width 79
func (this *TMessage) RplFwdString() (string) {
	var intermediate string = "NML -          -          - -"
	if this.ForwardFrom != nil { intermediate = fmt.Sprintf("FWD -          %s", this.ForwardFrom.PrintableString()) }
	if this.ForwardFromChat != nil { intermediate = fmt.Sprintf("FWD -          %s", this.ForwardFromChat.PrintableString()) }

	if this.ReplyToMessage != nil && this.ReplyToMessage.From != nil { intermediate = fmt.Sprintf("RPL %s %s", this.ReplyToMessage.IdString(), this.ReplyToMessage.From.PrintableString()) }
	if this.ReplyToMessage != nil && this.ReplyToMessage.From == nil { intermediate = fmt.Sprintf("RPL %s %s", this.ReplyToMessage.IdString(), this.ReplyToMessage.Chat.PrintableString()) } // longest, width 79

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
	} else if this.NewChatMembers != nil { intermediate = "ADDUSER"
	} else if this.LeftChatMember != nil { intermediate = "DELUSER"
	} else if this.NewChatTitle != nil { intermediate = "SETTITLE"
	} else if this.NewChatPhoto != nil { intermediate = "SETPHOTO"
	} else if this.DeleteChatPhoto != nil { intermediate = "DELPHOTO"
	} else if this.GroupChatCreated != nil { intermediate = "MKGROUP"
	} else if this.SupergroupChatCreated != nil { intermediate = "MKSGROUP"
	} else if this.ChannelChatCreated != nil { intermediate = "MKCHANNL"
	} else if this.MigrateToChatId != nil { intermediate = "TOSGRP"
	} else if this.MigrateFromChatId != nil { intermediate = "FROMSGRP"
	} else if this.PinnedMessage != nil { intermediate = "PINMSG" }
	return fmt.Sprintf("%-8s", intermediate)
}

// width: 41
func (this *TMessage) TimestampString() (string) {
	tf_string := "2006-01-02 15:04:05"
	if this.ForwardDate != nil {
		return fmt.Sprintf("%s (%s)", time.Unix(this.Date, 0).UTC().Format(tf_string), time.Unix(*this.ForwardDate, 0).UTC().Format(tf_string))
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

func userStrings(users []TUser) string {
	var names []string
	for _, u := range users {
		names = append(names, u.PrintableString())
	}
	return strings.Join(names, " ")
}

// width: variable.
func (this *TMessage) MessageContents() (string) {
	var caption string = "(no caption)"
	if this.Caption != nil { caption = fmt.Sprintf("<%s>", strings.Replace(*this.Caption, "\n", "", -1)) }

	if this.Audio != nil { return string(this.Audio.Id) }
	if this.Document != nil { return fmt.Sprintf("%s %s", this.Document.Id, caption) }
	if this.Game != nil { return this.Game.Title }
	if this.Photo != nil { return fmt.Sprintf("%s %s", GetLargestPhoto(this.Photo).Id, caption) }
	if this.Sticker != nil { return string(this.Sticker.Id) }
	if this.Video != nil { return fmt.Sprintf("%s %s", this.Video.Id, caption) }
	if this.Voice != nil { return string(this.Voice.Id) }
	if this.Contact != nil { return this.Contact.PhoneNumber + " " + this.Contact.FirstName }
	if this.Location != nil { return fmt.Sprintf("(%f %f)", this.Location.Longitude, this.Location.Latitude) }
	if this.Venue != nil { return fmt.Sprintf("(%f %f) %s: %s", this.Venue.Location.Longitude, this.Venue.Location.Latitude, this.Venue.Title, this.Venue.Description) }
	if this.NewChatMembers != nil { return fmt.Sprintf("[Chat member added] %s", userStrings(*this.NewChatMembers)) }
	if this.LeftChatMember != nil { return fmt.Sprintf("[Chat member removed] %s", this.LeftChatMember.PrintableString()) }
	if this.NewChatTitle != nil { return fmt.Sprintf("[New chat title] %s", *this.NewChatTitle) }
	if this.NewChatPhoto != nil { return fmt.Sprintf("[Chat photo set] %s", GetLargestPhoto(this.NewChatPhoto).Id) }
	if this.DeleteChatPhoto != nil { return "[Chat photo deleted]" }
	if this.GroupChatCreated != nil { return "[Group created]" }
	if this.SupergroupChatCreated != nil { return "[Supergroup created]" }
	if this.ChannelChatCreated != nil { return "[Channel created]" }
	if this.MigrateToChatId != nil { return fmt.Sprintf("[supergroup created from existing chat] %d", this.MigrateToChatId) }
	if this.MigrateFromChatId != nil { return fmt.Sprintf("[Group converted to supergroup] %d", this.MigrateFromChatId) }
	if this.PinnedMessage != nil { return fmt.Sprintf("[Message pinned] %s", this.PinnedMessage.MessageContents()) }
	if this.Text != nil { return strings.Replace(*this.Text, "\n", "", -1) }
	return "Unknown"
}

func (this *TMessage) PlainText() (string) {
	if this.Text != nil {
		return *this.Text
	} else if this.Caption != nil {
		return *this.Caption
	} else {
		return ""
	}
}

func (this *TMessage) Sender() (Sender) {
	if this.From != nil {
		return Sender{User: this.From.Id}
	} else {
		return Sender{Channel: this.Chat.Id}
	}
}

func (this Sender) String() (string) {
	if this.User != UserID(0) {
		return fmt.Sprintf("User[%d]", this.User)
	} else {
		return fmt.Sprintf("Channel[%d]", this.Channel)
	}
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

func (this *ORestrict) ToTChatPermissions() (TChatPermissions) {
	return TChatPermissions{
		CanSendMessages: this.CanSendMessages,
		CanSendMedia: this.CanSendMedia,
		CanSendPolls: this.CanSendPolls,
		CanSendOther: this.CanSendInline,
		CanPreviewLinks: this.CanSendWebPreviews,
		CanChangeInfo: this.CanChangeInfo,
		CanInviteUsers: this.CanInviteUsers,
		CanPinMessages: this.CanPinMessages,
	}
}

func (this *ORestrict) FromTChatPermissions(permissions TChatPermissions, chat_id int64) {
	*this = ORestrict{
		ChatID: chat_id,
		CanSendMessages: permissions.CanSendMessages,
		CanSendMedia: permissions.CanSendMedia,
		CanSendPolls: permissions.CanSendPolls,
		CanSendInline: permissions.CanSendOther,
		CanSendWebPreviews: permissions.CanPreviewLinks,
		CanChangeInfo: permissions.CanChangeInfo,
		CanInviteUsers: permissions.CanInviteUsers,
		CanPinMessages: permissions.CanPinMessages,
	}
}
