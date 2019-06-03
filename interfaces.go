package gogram


type InlineQueryable interface {
	ProcessInlineQuery(*InlineCtx)
	ProcessInlineQueryResult(*InlineResultCtx)
}


type Callbackable interface {
	ProcessCallback(*CallbackCtx)
}


type Messagable interface {
	ProcessMessage(*MessageCtx)
}


type Maintainer interface {
	DoMaintenance(*TelegramBot)
	GetInterval() int64
}


type InitSettings interface {
	InitializeAll(*TelegramBot) error
}
