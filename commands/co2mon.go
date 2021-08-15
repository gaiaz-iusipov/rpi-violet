package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/gaiaz-iusipov/rpi-violet/internal/monitor"
	"github.com/gaiaz-iusipov/rpi-violet/internal/monitor/co2mon"
)

var co2monCmd = &cobra.Command{
	Use: "co2mon",
	RunE: func(cmd *cobra.Command, args []string) error {
		dev, err := co2mon.Open(co2mon.WithRandomKey())
		if err != nil {
			return fmt.Errorf("co2mon.Open(): %w", err)
		}
		defer dev.Close()

		mon := monitor.New(dev)
		defer mon.Close()

		done := make(chan struct{})
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					if m, ok := mon.Measurements(); ok {
						fmt.Println(m)
					}
				}
			}
		}()

		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
		<-termChan

		close(done)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(co2monCmd)
}
