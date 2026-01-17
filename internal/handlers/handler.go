package handlers

import (
	"encoding/json"
	"go-simple-tg-bot/internal/clients"
	"go-simple-tg-bot/internal/models"
	"go-simple-tg-bot/internal/utils"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

type TelegramClient interface {
	Updates(offset, limit int) ([]models.Update, error)
	SendMessage(chatID int, text string) error
	SendPhotoByURL(chatID int, photoURL, caption string) error
}

type Handler struct {
	bot *clients.Client
	log *slog.Logger
}

func New(bot *clients.Client, log *slog.Logger) *Handler {
	return &Handler{
		bot: bot,
		log: log,
	}
}

func (h *Handler) HandleUpdate(update models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	if strings.HasPrefix(text, "/start") || strings.HasPrefix(text, "/help") {
		h.bot.SendMessage(chatID, `
Простой телеграм бот

/start - начать работу с ботом
/current - текущие показания температуры и влажности (не доступна)
/echo <ваше_сообщение> - повторяет за вами 
/dog - отправляет случайную фотографию с собакой
/help - справка по командам
		`)
		return
	}

	if strings.HasPrefix(text, "/echo") {
		msg := strings.TrimPrefix(text, "/echo ")
		h.bot.SendMessage(chatID, "Эхо: "+msg)
		return
	}

	if strings.HasPrefix(text, "/dog") {
		h.sendDog(chatID)
		return
	}

	h.bot.SendMessage(chatID, "Неизвестная команда. Используйте /help")
}

func (h *Handler) sendDog(chatID int) {
	req, err := http.NewRequest(http.MethodGet, "https://random.dog/woof.json", nil)
	if err != nil {
		h.log.Info("Ошибка создания запроса", utils.Err(err))
		h.bot.SendMessage(chatID, "Не удалось создать запрос. Попробуйте позже.")
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.log.Info("Ошибка при получении ответа", utils.Err(err))
		h.bot.SendMessage(chatID, "Не удалось получить фотографию собаки. Попробуйте позже.")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.log.Info("Ошибка при чтении тела ответа", utils.Err(err))
		h.bot.SendMessage(chatID, "Ошибка при обработке данных. Попробуйте позже.")
		return
	}

	type randomDogResponse struct {
		URL string `json:"url"`
	}

	var dogData randomDogResponse
	err = json.Unmarshal(body, &dogData)
	if err != nil {
		h.log.Info("Не удалось обработать ответ", utils.Err(err))
		h.bot.SendMessage(chatID, "Ошибка формата данных. Попробуйте позже.")
		return
	}
	log.Println(dogData.URL)

	err = h.bot.SendPhotoByURL(chatID, dogData.URL, "")
	if err != nil {
		h.log.Info("Ошибка при отправке сообщения пользователю", slog.Int("chat_id", chatID), utils.Err(err))
		h.bot.SendMessage(chatID, "Ошибка при отправке фотографии. Попробуйте позже.")
	}
}
