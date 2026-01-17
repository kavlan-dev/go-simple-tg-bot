# Telegram Bot на Go

Простой Telegram бот, реализованный на языке Go с использованием Telegram Bot API.

## Описание

Этот проект представляет собой простого Telegram бота, который поддерживает следующие команды:

- `/start` - Начать работу с ботом
- `/help` - Показать справку по командам
- `/echo <сообщение>` - Повторить ваше сообщение
- `/dog` - Отправить случайную фотографию собаки

## Структура проекта

```
go-simple-tg-bot/
├── cmd/
│   └── bot/
│       └── main.go          # Основной файл запуска бота
├── internal/
│   ├── clients/
│   │   └── client.go        # Клиент для работы с Telegram API
│   ├── handlers/
│   │   └── handler.go       # Обработчики сообщений
│   ├── models/
│   │   └── update.go        # Модели данных
│   └── utils/
│       └── token.go         # Утилиты для работы с токенами
├── go.mod                   # Файл зависимостей
└── README.md                # Документация
```

## Установка и запуск

### Предварительные требования

- Go 1.25 или выше
- Telegram Bot Token (можно получить у [@BotFather](https://t.me/BotFather))

### Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/go-simple-tg-bot.git
   cd go-simple-tg-bot
   ```

2. Установите зависимости:
   ```bash
   go mod tidy
   ```

### Запуск

Есть два способа запуска бота:

#### Способ 1: Использование переменной окружения

```bash
export TOKEN=your_telegram_bot_token
go run cmd/bot/main.go
```

#### Способ 2: Использование флага командной строки

```bash
go run cmd/bot/main.go -t your_telegram_bot_token
```

## Архитектура

Проект следует принципам чистой архитектуры и разделения ответственности:

1. **Client** (`internal/clients/client.go`): Отвечает за взаимодействие с Telegram API
2. **Handler** (`internal/handlers/handler.go`): Обрабатывает входящие сообщения и команды
3. **Models** (`internal/models/update.go`): Содержит структуры данных для работы с API
4. **Utils** (`internal/utils/token.go`): Вспомогательные функции

## Команды бота

| Команда | Описание | Пример |
|---------|----------|--------|
| `/start` | Начать работу с ботом | `/start` |
| `/help` | Показать справку | `/help` |
| `/echo` | Повторить сообщение | `/echo Привет!` |
| `/dog` | Отправить случайную фотографию собаки | `/dog` |

## Лицензия

MIT License
