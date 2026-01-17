package utils

import (
	"flag"
	"log"
	"os"
)

func MustToken() string {
	token := os.Getenv("TOKEN")
	if token != "" {
		return token
	}

	token = *flag.String(
		"t",
		"",
		"Токен телеграм бота",
	)
	flag.Parse()

	if token == "" {
		log.Fatal("Токен не указан")
	}

	return token
}
