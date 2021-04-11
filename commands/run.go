package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
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

	location, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return fmt.Errorf("time.LoadLocation: %w", err)
	}

	c := cron.New(
		cron.WithLocation(location),
	)

	_, err = c.AddJob("@daily", &cronJob{
		pin:         pin,
		state:       gpio.Low,
		photoSender: tg,
	})
	if err != nil {
		return err
	}

	_, err = c.AddJob("0 7 * * *", &cronJob{
		pin:         pin,
		state:       gpio.High,
		photoSender: tg,
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

type PhotoSender interface {
	SendPhoto(context.Context, []byte) error
}

type cronJob struct {
	pin         gpio.PinIO
	state       gpio.Level
	photoSender PhotoSender
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

	pic, err := makePicture(ctx)
	if err != nil {
		return fmt.Errorf("failed to make photo: %w", err)
	}

	err = cj.photoSender.SendPhoto(ctx, pic)
	if err != nil {
		return fmt.Errorf("failed to send photo: %w", err)
	}
	return nil
}

func makePicture(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Raspistill.Timeout))
	defer cancel()

	strQuality := strconv.Itoa(int(cfg.Raspistill.JpegQuality))
	cmd := exec.CommandContext(ctx, "raspistill", "-q", strQuality, "-o", "-")

	out, err := cmd.Output()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, errors.New("timed out")
		}
		return nil, err
	}
	return out, nil
}
