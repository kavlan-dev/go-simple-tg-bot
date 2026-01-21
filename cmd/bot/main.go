package main

import (
	"context"
	"go-simple-tg-bot/internal/clients"
	"go-simple-tg-bot/internal/config"
	"go-simple-tg-bot/internal/handlers"
	"go-simple-tg-bot/internal/models"
	"go-simple-tg-bot/internal/utils"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	offset := 0

	go func() {
		<-sigChan
		log.Info("Получен сигнал завершения, начинаем плавное завершение...")
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		default:
			updates, err := tgClient.Updates(ctx, offset, 0)
			if err != nil {
				log.Info("Ошибка получения обновлений", utils.Err(err))
				time.Sleep(1 * time.Second)
				continue
			}

			for _, update := range updates {
				log.Info("Новый запрос", slog.Any("message", *update.Message))
				wg.Add(1)
				go func(u models.Update) {
					defer wg.Done()
					handler.HandleUpdate(ctx, u)
				}(update)
				offset = update.UpdateID + 1
				time.Sleep(1 * time.Second)
			}
		}
	}
}
