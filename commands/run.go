package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
	tb "gopkg.in/tucnak/telebot.v2"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run application",
	PreRunE: initConfig,
	RunE:    runE,
}

func init() {
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	rootCmd.AddCommand(runCmd)
}

func runE(_ *cobra.Command, _ []string) error {
	pin := rpio.Pin(cfg.GPIOLightPin)

	err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDNS,
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %w", err)
	}
	defer sentry.Flush(cfg.SentryTimeout)

	tgClient := &http.Client{
		Timeout: cfg.TelegramClientTimeout,
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  cfg.TelegramBotToken,
		Client: tgClient,
	})
	if err != nil {
		return fmt.Errorf("tb.NewBot: %w", err)
	}

	chat := &tb.Chat{ID: cfg.TelegramChatID}
	location, _ := time.LoadLocation("Europe/Moscow")
	c := cron.New(
		cron.WithLocation(location),
	)

	_, err = c.AddJob("@daily", &cronJob{
		pin:    pin,
		state:  rpio.Low,
		tgBot:  bot,
		tgChat: chat,
	})
	if err != nil {
		return err
	}

	_, err = c.AddJob("0 7 * * *", &cronJob{
		pin:    pin,
		state:  rpio.High,
		tgBot:  bot,
		tgChat: chat,
	})
	if err != nil {
		return err
	}

	c.Start()
	defer c.Stop()

	go func() {
		log.Println(http.ListenAndServe("localhost:"+cfg.DebugPort, nil))
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	return nil
}

type cronJob struct {
	pin    rpio.Pin
	state  rpio.State
	tgBot  *tb.Bot
	tgChat *tb.Chat
}

func (cj cronJob) Run() {
	if err := cj.runE(); err != nil {
		sentry.CaptureException(err)
	}
}

func (cj cronJob) runE() error {
	if err := rpio.Open(); err != nil {
		return fmt.Errorf("unable to open gpio: %w", err)
	}
	defer rpio.Close()

	cj.pin.Output()
	cj.pin.Write(cj.state)

	reader, err := makePicture(context.Background())
	if err != nil {
		return fmt.Errorf("failed to make photo: %w", err)
	}

	_, err = cj.tgBot.Send(cj.tgChat, &tb.Photo{File: tb.FromReader(reader)})
	if err != nil {
		return fmt.Errorf("bot.Send daily message: %w", err)
	}
	return nil
}

func makePicture(ctx context.Context) (*bytes.Reader, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.RaspistillTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "raspistill", "-o", "-")

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil, errors.New("timed out")
	}

	return bytes.NewReader(out), nil
}
