package bot

import (
	"github.com/BlackRRR/checker-bot/internal/app/model"
	"github.com/BlackRRR/checker-bot/internal/db/redis"
	"github.com/bots-empire/base-bot/msgs"
)

func (b *BotService) Admin(s *model.Situation) error {
	redis.RdbSetUser(s.User.ID, "/admin")
	//admins, err := b.getAdmins(&model.GetAdmins{
	//	Code: b.GlobalBot.BotReq,
	//})
	//if err != nil {
	//	return err
	//}

	//for _, val := range admins {
	//	if val == s.User.ID {
	//		markUp := msgs.NewIlMarkUp(
	//			msgs.NewIlRow(msgs.NewIlDataButton("set_url", "/url")),
	//			msgs.NewIlRow(msgs.NewIlDataButton("set_text", "/text")),
	//		).Build(b.GlobalBot.Language[s.BotLang])
	//
	//		text := b.GlobalBot.LangText(s.BotLang, "admin_settings")
	//		err := b.BaseBotSrv.NewParseMarkUpMessage(s.User.ID, &markUp, text)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//}

	if 872383555 == s.User.ID {
		markUp := msgs.NewIlMarkUp(
			msgs.NewIlRow(msgs.NewIlDataButton("set_url", "/url")),
			msgs.NewIlRow(msgs.NewIlDataButton("set_text", "/text")),
		).Build(b.GlobalBot.Language[s.BotLang])

		text := b.GlobalBot.LangText(s.BotLang, "admin_settings")
		err := b.BaseBotSrv.NewParseMarkUpMessage(s.User.ID, &markUp, text)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BotService) SetText(s *model.Situation) error {
	text, err := b.Repo.GetText()
	if err != nil {
		return err
	}

	redis.RdbSetUser(s.User.ID, "/set_text")

	if text == "" {
		text := b.GlobalBot.LangText(s.BotLang, "no_text")
		return b.BaseBotSrv.NewParseMessage(s.User.ID, text)
	}

	text = b.GlobalBot.LangText(s.BotLang, "new_text", text)
	return b.BaseBotSrv.NewParseMessage(s.User.ID, text)
}

func (b *BotService) TextReady(s *model.Situation) error {
	text, err := b.Repo.GetText()
	if err != nil {
		return err
	}

	if text == "" {
		err := b.Repo.SetText(s.Message.Text)
		if err != nil {
			return err
		}
	}

	err = b.Repo.UpdateText(s.Message.Text)
	if err != nil {
		return err
	}

	text = b.GlobalBot.LangText(s.BotLang, "text_ready")
	return b.BaseBotSrv.NewParseMessage(s.User.ID, text)
}

func (b *BotService) SetURL(s *model.Situation) error {
	url, err := b.Repo.GetURL()
	if err != nil {
		return err
	}

	redis.RdbSetUser(s.User.ID, "/set_url")

	if url == "" {
		text := b.GlobalBot.LangText(s.BotLang, "no_url")
		return b.BaseBotSrv.NewParseMessage(s.User.ID, text)
	}

	text := b.GlobalBot.LangText(s.BotLang, "new_url", url)
	return b.BaseBotSrv.NewParseMessage(s.User.ID, text)
}

func (b *BotService) URLReady(s *model.Situation) error {
	url, err := b.Repo.GetURL()
	if err != nil {
		return err
	}

	if url == "" {
		err := b.Repo.SetURL(s.Message.Text)
		if err != nil {
			return err
		}
	}

	err = b.Repo.UpdateURL(s.Message.Text)
	if err != nil {
		return err
	}

	text := b.GlobalBot.LangText(s.BotLang, "url_ready")
	return b.BaseBotSrv.NewParseMessage(s.User.ID, text)
}
