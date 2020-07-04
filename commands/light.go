package commands

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
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

func lightE(action lightAction) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := rpio.Open(); err != nil {
			return fmt.Errorf("unable to open gpio: %w", err)
		}
		defer rpio.Close()

		pin := rpio.Pin(cfg.GPIOLightPin)
		pin.Output()

		switch action {
		case actionOn:
			pin.High()
		case actionOff:
			pin.Low()
		case actionToggle:
			pin.Toggle()
		}
		return nil
	}
}
