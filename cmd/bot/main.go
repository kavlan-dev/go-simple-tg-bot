package main

import (
	"go-simple-tg-bot/internal/clients"
	"go-simple-tg-bot/internal/config"
	"go-simple-tg-bot/internal/handlers"
	"go-simple-tg-bot/internal/utils"
	"log"
	"log/slog"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln("Не удалось загрузить файл конфигураций", err)
	}

	tgClient := clients.New("api.telegram.org", cfg.Token)
	log := utils.InitLogger(cfg.Env)

	handler := handlers.New(tgClient, log)

	offset := 0

	for {
		updates, err := tgClient.Updates(offset, 0)
		if err != nil {
			log.Info("Ошибка получения обновлений: %v", utils.Err(err))
		}
		for _, update := range updates {
			log.Info("Новый запрос", slog.Any("message", *update.Message))
			go handler.HandleUpdate(update)
			offset = update.UpdateID + 1
			time.Sleep(1 * time.Second)
		}
	}
}
