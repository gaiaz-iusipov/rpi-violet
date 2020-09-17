package commands

import (
	"errors"
	"fmt"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

type lightAction int

const (
	actionOn lightAction = iota
	actionOff
	actionToggle
)

func (a lightAction) String() string {
	return [...]string{"on", "off", "toggle"}[a]
}

var (
	lightCmd = &cobra.Command{
		Use:               "light",
		Short:             "Light switch",
		PersistentPreRunE: initConfig,
	}
	lightOnCmd = &cobra.Command{
		Use:  actionOn.String(),
		RunE: lightE(actionOn),
	}
	lightOffCmd = &cobra.Command{
		Use:  actionOff.String(),
		RunE: lightE(actionOff),
	}
	lightToggleCmd = &cobra.Command{
		Use:  actionToggle.String(),
		RunE: lightE(actionToggle),
	}
)

func init() {
	lightCmd.AddCommand(lightOnCmd, lightOffCmd, lightToggleCmd)
	lightCmd.Flags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	rootCmd.AddCommand(lightCmd)
}

func lightE(action lightAction) func(_ *cobra.Command, _ []string) error {
	return func(_ *cobra.Command, _ []string) error {
		if _, err := host.Init(); err != nil {
			return fmt.Errorf("periph.Init: %w", err)
		}

		pin := gpioreg.ByName(cfg.GPIOLightPin)
		if pin == nil {
			return errors.New("gpio pin is not present")
		}

		var lvl gpio.Level

		switch action {
		case actionOn:
			lvl = gpio.High
		case actionOff:
			lvl = gpio.Low
		case actionToggle:
			lvl = !pin.Read()
		}

		if err := pin.Out(lvl); err != nil {
			return fmt.Errorf("pin.Out: %w", err)
		}

		return nil
	}
}
