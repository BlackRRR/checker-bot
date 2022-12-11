package layers

import (
	"github.com/BlackRRR/checker-bot/internal/app/model"
	"github.com/BlackRRR/checker-bot/internal/app/services/bot"
)

type AdminHandlers struct {
	Handlers   map[string]model.Handler
	BotService *bot.BotService
}

func (h *AdminHandlers) GetHandler(command string) model.Handler {
	return h.Handlers[command]
}

func (h *AdminHandlers) Init() {
	//admin message
}

func (h *AdminHandlers) OnCommand(command string, handler model.Handler) {
	h.Handlers[command] = handler
}
