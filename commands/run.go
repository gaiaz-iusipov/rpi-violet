package commands

import (
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
	"github.com/spf13/cobra"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"

	"github.com/gaiaz-iusipov/rpi-violet/internal/config"
	"github.com/gaiaz-iusipov/rpi-violet/internal/cron"
	"github.com/gaiaz-iusipov/rpi-violet/internal/raspistill"
	"github.com/gaiaz-iusipov/rpi-violet/internal/telegram"
	"github.com/gaiaz-iusipov/rpi-violet/internal/version"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	RunE:  runE,
}

var cfgFile string

func init() {
	runCmd.Flags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
	rootCmd.AddCommand(runCmd)
}

func runE(_ *cobra.Command, _ []string) error {
	cfg, err := config.New(cfgFile)
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

	c, err := cron.New(cfg.Cron, location, pin, rs, tg)
	if err != nil {
		return fmt.Errorf("cron.New: %w", err)
	}
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
