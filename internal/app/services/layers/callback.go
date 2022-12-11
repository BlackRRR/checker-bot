package layers

import (
	"github.com/BlackRRR/checker-bot/internal/app/model"
	"github.com/BlackRRR/checker-bot/internal/app/services/bot"
)

type CallBackHandlers struct {
	Handlers   map[string]model.Handler
	BotService *bot.BotService
}

func (h *CallBackHandlers) GetHandler(command string) model.Handler {
	return h.Handlers[command]
}

func (h *CallBackHandlers) Init() {
	//Money command
	h.OnCommand("/url", h.BotService.SetURL)
	h.OnCommand("/text", h.BotService.SetText)

}

func (h *CallBackHandlers) OnCommand(command string, handler model.Handler) {
	h.Handlers[command] = handler
}
