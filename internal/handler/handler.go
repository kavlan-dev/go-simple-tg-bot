package handler

import (
	"context"
	"go-simple-tg-bot/internal/model"
	"go-simple-tg-bot/internal/util"
	"log/slog"
	"strings"
)

type telegramClient interface {
	Updates(ctx context.Context, offset, limit int) ([]model.Update, error)
	SendMessage(ctx context.Context, chatID int, text string) error
	SendPhotoByURL(ctx context.Context, chatID int, photoURL, caption string) error
}

type service interface {
	DogImage(ctx context.Context) (string, error)
}

type handler struct {
	bot     telegramClient
	service service
	log     *slog.Logger
}

func NewHandler(bot telegramClient, service service, log *slog.Logger) *handler {
	return &handler{
		bot:     bot,
		service: service,
		log:     log,
	}
}

func (h *handler) HandleUpdate(ctx context.Context, update model.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	if strings.HasPrefix(text, "/start") || strings.HasPrefix(text, "/help") {
		h.bot.SendMessage(ctx, chatID, `
Простой телеграм бот

/start - начать работу с ботом
/dog - отправляет случайную фотографию с собакой
/help - справка по командам
		`)
		return
	}

	if strings.HasPrefix(text, "/dog") {
		h.sendDog(ctx, chatID)
		return
	}

	h.bot.SendMessage(ctx, chatID, "Неизвестная команда. Используйте /help")
}

func (h *handler) sendDog(ctx context.Context, chatID int) {
	url, err := h.service.DogImage(ctx)
	if err != nil {
		h.log.Info("Ошибка при получении ссылки с фотографией", util.Err(err))
		h.bot.SendMessage(ctx, chatID, "Неизвестная команда. Используйте /help")
	}

	if err := h.bot.SendPhotoByURL(ctx, chatID, url, ""); err != nil {
		h.log.Info("Ошибка при отправке сообщения пользователю", slog.Int("chat_id", chatID), util.Err(err))
		h.bot.SendMessage(ctx, chatID, "Ошибка при отправке фотографии. Попробуйте позже.")
	}
}
