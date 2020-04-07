package main

// "github.com/br9k777/telnet/pkg/telnet"
import (
	// "database/sql"
	// "errors"

	"fmt"
	"os"

	"time"

	"github.com/br9k777/telnet/pkg/config"
	"github.com/br9k777/telnet/pkg/server"
	"github.com/br9k777/telnet/pkg/telnet"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {

	logger, err := config.GetStandartLogger("production")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Can't create logger %s", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)

	var timeout time.Duration
	app := &cli.App{
		Name:      "telnet",
		Version:   "v0.1.0",
		Usage:     "simple telnet",
		ArgsUsage: "host port",
		Authors: []*cli.Author{
			{Name: "br9k"},
		},
	}

	app.Flags = []cli.Flag{
		&cli.DurationFlag{
			Name:        "timeout",
			Aliases:     []string{"w"},
			Value:       10 * time.Second,
			Usage:       "`DURATION` of the attempt to connect to the server",
			Destination: &timeout,
			EnvVars:     []string{"TELNET_TIMEOUT", "TIMEOUT"},
		},
		&cli.BoolFlag{
			Name:    "work-as-server",
			Aliases: []string{"server", "s"},
			Value:   false,
			Usage:   "work as repeat server for tests",
			EnvVars: []string{"TELNET_TIMEOUT", "TIMEOUT"},
		},
	}
	app.Action = func(c *cli.Context) (err error) {

		if c.Bool("work-as-server") {
			err = server.StartServer(c.Args().Get(0), c.Args().Get(1))
		} else {
			zap.S().Infof("Timeout =%s\n", timeout.String())
			err = telnet.StartTelnetClient(timeout, c.Args().Get(0), c.Args().Get(1))
		}
		return err
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s\n", err)
	}
}
