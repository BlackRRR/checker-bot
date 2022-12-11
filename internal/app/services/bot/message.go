package bot

import (
	"bytes"
	"encoding/json"
	"github.com/BlackRRR/checker-bot/internal/app/model"
	"github.com/bots-empire/base-bot/msgs"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

func (b *BotService) StartCommand(s *model.Situation) error {
	info, err := b.getIncomeInfo(s.Message.From.ID)
	if err != nil {
		return err
	}

	var text string
	if info == nil {
		text = b.GlobalBot.LangText(s.BotLang, "user_info_unknown_admin",
			s.Message.From.ID,
			s.Message.From.UserName,
		)
	} else {
		text = b.GlobalBot.LangText(s.BotLang, "user_info_admin",
			s.Message.From.ID,
			s.Message.From.UserName,
			info.IncomeSource,
			info.BotName,
			info.TypeBot,
			info.BotLink)
	}

	err = b.Repo.SaveIncomeInfo(&model.IncomeInfo{
		UserID:       s.Message.From.ID,
		BotLink:      info.BotLink,
		BotName:      info.BotName,
		IncomeSource: info.IncomeSource,
		TypeBot:      info.TypeBot,
	}, s.Message.From.UserName)
	if err != nil {
		return err
	}

	req := &model.GetAdmins{
		Code: b.GlobalBot.BotReq,
	}

	admins, err := b.getAdmins(req)
	if err != nil {
		return err
	}

	duplicates := make(map[int64]struct{}, 0)
	for _, val := range admins {
		if _, ok := duplicates[val.UserID]; !ok {
			duplicates[val.UserID] = struct{}{}
		} else {
			continue
		}

		err = b.BaseBotSrv.NewParseMessage(val.UserID, text)
		if err != nil {
			return err
		}
	}

	//text = b.GlobalBot.LangText(s.BotLang, "starting")
	url, err := b.Repo.GetURL()
	if err != nil {
		return err
	}

	if url == "" {
		return nil
	}

	text, err = b.Repo.GetText()
	if err != nil {
		return err
	}

	markUp := msgs.NewIlMarkUp(msgs.NewIlRow(msgs.NewIlURLButton("go_to_url", url))).Build(b.GlobalBot.Language[b.GlobalBot.BotLang])

	return b.BaseBotSrv.NewParseMarkUpMessage(s.User.ID, &markUp, text)

	//for range time.Tick(time.Millisecond * 100) {
	//	progressBar := 0
	//	a := 10
	//	text = b.GlobalBot.LangText(s.BotLang, "time_text", progressBar+a+rand.Intn(10))
	//	a *= 2
	//
	//	if a >= 80 {
	//		text = b.GlobalBot.LangText(s.BotLang, "close_point")
	//		err = b.BaseBotSrv.NewParseMessage(s.User.ID, text)
	//		if err != nil {
	//			return err
	//		}
	//
	//		time.Sleep(time.Millisecond * 700)
	//	}
	//
	//	if a >= 100 {
	//		progressBar = 100
	//		err = b.BaseBotSrv.NewParseMessage(s.User.ID, text)
	//		if err != nil {
	//			return err
	//		}
	//
	//		text = b.GlobalBot.LangText(s.BotLang, "success")
	//		err = b.BaseBotSrv.NewParseMessage(s.User.ID, text)
	//		if err != nil {
	//			return err
	//		}
	//
	//		url, err := b.Repo.GetURL()
	//		if err != nil {
	//			return err
	//		}
	//
	//		if url == "" {
	//			return nil
	//		}
	//
	//		text = b.GlobalBot.LangText(s.BotLang, "your_url", url)
	//		return b.BaseBotSrv.NewParseMessage(s.User.ID, text)
	//	}

}

func (b *BotService) getIncomeInfo(userID int64) (*model.IncomeInfo, error) {
	req := &model.GetIncomeInfo{
		UserID:  userID,
		TypeBot: b.GlobalBot.BotReq,
	}

	marshal, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(marshal)

	resp, err := http.Post(GoHttp+b.Config.Server.IP+b.Config.Server.IncomeInfoRoute, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Fatalf("Get Income info: Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, data)
	}

	var info *model.IncomeInfo

	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal income info")
	}

	return info, nil
}

func (b *BotService) getAdmins(req *model.GetAdmins) ([]*model.Access, error) {
	//marshal, err := json.Marshal(req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//body := bytes.NewReader(marshal)

	resp, err := http.Post(GoHttp+b.Config.Server.IP+b.Config.Server.AdminRout, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Fatalf("Get admins: Response failed with status code: %d and\nreq: %s\n", resp.StatusCode, data)
	}

	var adminIDs []*model.Access

	err = json.Unmarshal(data, &adminIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal admin ids")
	}

	return adminIDs, nil
}
