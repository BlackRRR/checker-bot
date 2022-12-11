package main

import (
	"github.com/BlackRRR/checker-bot/config"
	"github.com/BlackRRR/checker-bot/internal/app/model"
	"github.com/BlackRRR/checker-bot/internal/app/repository"
	"github.com/BlackRRR/checker-bot/internal/app/services"
	"github.com/BlackRRR/checker-bot/internal/app/services/bot"
	"github.com/BlackRRR/checker-bot/internal/app/utils"
	"github.com/BlackRRR/checker-bot/internal/db"
	"github.com/BlackRRR/checker-bot/internal/log"
	"github.com/bots-empire/base-bot/msgs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log2 "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//init logger
	logger := log.NewDefaultLogger().Prefix("Checker Bot")
	log.PrintLogo("Checker Bot", []string{"8000FF"})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		metricErr := http.ListenAndServe(":7041", nil)
		if metricErr != nil {
			logger.Fatal("metrics stopped by metricErr: %s", metricErr.Error())
		}
	}()

	//Init Config
	cfg, dbConn, err := config.InitConfig()
	if err != nil {
		log2.Fatal(err)
	}

	//Init Database
	pool, err := db.InitDB(cfg.Pgx, dbConn)
	if err != nil {
		log2.Fatal(err)
	}

	//init bots config
	srvs := make([]*services.Services, 0)
	for _, bots := range cfg.Bots {
		for _, language := range bots.Components {
			globalBot := model.FillBotsConfig(language.Token, language.BotLang, bots.AMSBotType)
			msgsSrv := msgs.NewService(globalBot, []int64{872383555, 1418862576})
			repo := repository.NewRepository(pool, msgsSrv, globalBot)
			initServices := services.InitServices(repo, msgsSrv, globalBot, cfg)
			srvs = append(srvs, initServices)

		}
	}

	for _, service := range srvs {
		go func(handler *bot.BotService) {
			handler.ActionsWithUpdates(logger, utils.NewSpreader(time.Minute))
		}(service.BotSrv)

		service.BotSrv.BaseBotSrv.SendNotificationToDeveloper("Bot are restart", false)

		logger.Ok("service are running")
	}

	sig := <-subscribeToSystemSignals()

	log2.Printf("shutdown all process on '%s' system signal\n", sig.String())
}

func subscribeToSystemSignals() chan os.Signal {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	return ch
}
