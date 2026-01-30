package handlers

import (
	"context"
	"encoding/json"
	"go-simple-tg-bot/internal/models"
	"go-simple-tg-bot/internal/utils"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type telegramClient interface {
	Updates(ctx context.Context, offset, limit int) ([]models.Update, error)
	SendMessage(ctx context.Context, chatID int, text string) error
	SendPhotoByURL(ctx context.Context, chatID int, photoURL, caption string) error
}

type handler struct {
	bot telegramClient
	log *slog.Logger
}

func New(bot telegramClient, log *slog.Logger) *handler {
	return &handler{
		bot: bot,
		log: log,
	}
}

func (h *handler) HandleUpdate(ctx context.Context, update models.Update) {
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://random.dog/woof.json", nil)
	if err != nil {
		h.log.Info("Ошибка создания запроса", utils.Err(err))
		h.bot.SendMessage(ctx, chatID, "Не удалось создать запрос. Попробуйте позже.")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		h.log.Info("Ошибка при получении ответа", utils.Err(err))
		h.bot.SendMessage(ctx, chatID, "Не удалось получить фотографию собаки. Попробуйте позже.")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.log.Info("Неожиданный статус код от API", slog.Int("status_code", resp.StatusCode))
		h.bot.SendMessage(ctx, chatID, "Сервис временно недоступен. Попробуйте позже.")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.log.Info("Ошибка при чтении тела ответа", utils.Err(err))
		h.bot.SendMessage(ctx, chatID, "Ошибка при обработке данных. Попробуйте позже.")
		return
	}

	type randomDogResponse struct {
		URL string `json:"url"`
	}

	var dogData randomDogResponse
	err = json.Unmarshal(body, &dogData)
	if err != nil {
		h.log.Info("Не удалось обработать ответ", utils.Err(err))
		h.bot.SendMessage(ctx, chatID, "Ошибка формата данных. Попробуйте позже.")
		return
	}

	if dogData.URL == "" {
		h.log.Info("Получен пустой URL фотографии")
		h.bot.SendMessage(ctx, chatID, "Не удалось получить фотографию собаки. Попробуйте позже.")
		return
	}

	err = h.bot.SendPhotoByURL(ctx, chatID, dogData.URL, "")
	if err != nil {
		h.log.Info("Ошибка при отправке сообщения пользователю", slog.Int("chat_id", chatID), utils.Err(err))
		h.bot.SendMessage(ctx, chatID, "Ошибка при отправке фотографии. Попробуйте позже.")
	}
}
