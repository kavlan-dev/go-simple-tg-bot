package app

import (
	"context"
	"go-simple-tg-bot/internal/client"
	"go-simple-tg-bot/internal/config"
	"go-simple-tg-bot/internal/handler"
	"go-simple-tg-bot/internal/model"
	"go-simple-tg-bot/internal/service"
	"go-simple-tg-bot/internal/util"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalln("Не удалось загрузить файл конфигураций", err)
	}

	tgClient := client.NewClient("api.telegram.org", cfg.Token)
	if tgClient == nil {
		log.Fatalln("Не удалось создать Telegram клиент")
	}

	logger := util.InitLogger(cfg.Env)
	service := service.NewService()
	handler := handler.NewHandler(tgClient, service, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	offset := 0

	go func() {
		<-sigChan
		logger.Info("Получен сигнал завершения, начинаем плавное завершение...")
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
				logger.Error("Ошибка получения обновлений", util.Err(err))
				time.Sleep(1 * time.Second)
				continue
			}

			for _, update := range updates {
				logger.Info("Новый запрос", slog.Any("message", *update.Message))
				wg.Add(1)
				go func(u model.Update) {
					defer wg.Done()
					handler.HandleUpdate(ctx, u)
				}(update)
				offset = update.UpdateID + 1
				time.Sleep(1 * time.Second)
			}
		}
	}
}
