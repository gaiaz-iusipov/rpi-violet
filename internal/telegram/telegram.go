package telegram

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
)

type Telegram struct {
	bot  *tb.Bot
	chat *tb.Chat
}

func New(cfg *config.Telegram) (*Telegram, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.ClientTimeout),
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Client: httpClient,
	})
	if err != nil {
		return nil, fmt.Errorf("tb.NewBot: %w", err)
	}

	chat := &tb.Chat{ID: cfg.ChatID}
	return &Telegram{
		bot:  bot,
		chat: chat,
	}, nil
}

func (t *Telegram) SendPhoto(_ context.Context, photo []byte, caption string) error {
	return retry(10, 10*time.Second, func() error {
		thPhoto := &tb.Photo{
			File:    tb.FromReader(bytes.NewReader(photo)),
			Caption: caption,
		}

		_, err := t.bot.Send(t.chat, thPhoto)
		if err != nil {
			return fmt.Errorf("bot.Send: %w", err)
		}
		return nil
	})
}
