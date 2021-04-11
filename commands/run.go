package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
	"github.com/gaiaz-iusipov/rpi-violet/internal/raspistill"
	"github.com/gaiaz-iusipov/rpi-violet/internal/telegram"
	"github.com/gaiaz-iusipov/rpi-violet/pkg/version"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	RunE:  runE,
}

func init() {
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
	rootCmd.AddCommand(runCmd)
}

func runE(_ *cobra.Command, _ []string) error {
	var err error
	cfg, err = config.New(cfgFile)
	if err != nil {
		return fmt.Errorf("config.New: %w", err)
	}

	if _, err := host.Init(); err != nil {
		return fmt.Errorf("periph.Init: %w", err)
	}

	pin := gpioreg.ByName(cfg.GPIO.LightPin)
	if pin == nil {
		return errors.New("gpio pin is not present")
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:     cfg.Sentry.DNS,
		Release: "rpi-violet@" + version.Version(),
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %w", err)
	}
	defer sentry.Flush(time.Duration(cfg.Sentry.Timeout))

	tg, err := telegram.New(cfg.Telegram)
	if err != nil {
		return fmt.Errorf("bot.New: %w", err)
	}

	rs := raspistill.New(cfg.Raspistill)

	location, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return fmt.Errorf("time.LoadLocation: %w", err)
	}

	c := cron.New(
		cron.WithLocation(location),
	)

	_, err = c.AddJob("@daily", &cronJob{
		pin:           pin,
		state:         gpio.Low,
		photoProvider: rs,
		photoSender:   tg,
	})
	if err != nil {
		return err
	}

	_, err = c.AddJob("0 7 * * *", &cronJob{
		pin:           pin,
		state:         gpio.High,
		photoProvider: rs,
		photoSender:   tg,
	})
	if err != nil {
		return err
	}

	c.Start()
	defer c.Stop()

	go func() {
		debugPort := strconv.Itoa(int(cfg.DebugPort))
		log.Println(http.ListenAndServe("localhost:"+debugPort, nil))
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	return nil
}

type PhotoProvider interface {
	GetPhoto(context.Context) ([]byte, error)
}

type PhotoSender interface {
	SendPhoto(context.Context, []byte) error
}

type cronJob struct {
	pin           gpio.PinIO
	state         gpio.Level
	photoProvider PhotoProvider
	photoSender   PhotoSender
}

func (cj cronJob) Run() {
	if err := cj.runE(); err != nil {
		localHub := sentry.CurrentHub().Clone()
		localHub.CaptureException(err)
	}
}

func (cj cronJob) runE() error {
	if err := cj.pin.Out(cj.state); err != nil {
		return fmt.Errorf("pin.Out: %w", err)
	}

	ctx := context.Background()

	photo, err := cj.photoProvider.GetPhoto(ctx)
	if err != nil {
		return fmt.Errorf("failed to get photo: %w", err)
	}

	err = cj.photoSender.SendPhoto(ctx, photo)
	if err != nil {
		return fmt.Errorf("failed to send photo: %w", err)
	}
	return nil
}
