package layers

import (
	"github.com/BlackRRR/checker-bot/internal/app/model"
	"github.com/BlackRRR/checker-bot/internal/app/services/bot"
)

type MessagesHandlers struct {
	Handlers   map[string]model.Handler
	BotService *bot.BotService
}

func (h *MessagesHandlers) GetHandler(command string) model.Handler {
	return h.Handlers[command]
}

func (h *MessagesHandlers) Init() {
	//Start command
	h.OnCommand("/start", h.BotService.StartCommand)
	h.OnCommand("/admin", h.BotService.Admin)
	h.OnCommand("/set_url", h.BotService.URLReady)
	h.OnCommand("/set_text", h.BotService.TextReady)
}

func (h *MessagesHandlers) OnCommand(command string, handler model.Handler) {
	h.Handlers[command] = handler
}
