package command

import (
	"errors"

	"github.com/go-kit/kit/log/level"
	"github.com/promhippie/prometheus-vcd-sd/pkg/action"
	"github.com/promhippie/prometheus-vcd-sd/pkg/config"
	"github.com/urfave/cli/v2"
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "web.address",
				Value:       "0.0.0.0:9000",
				Usage:       "Address to bind the metrics server",
				EnvVars:     []string{"PROMETHEUS_VCD_WEB_ADDRESS"},
				Destination: &cfg.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "web.path",
				Value:       "/metrics",
				Usage:       "Path to bind the metrics server",
				EnvVars:     []string{"PROMETHEUS_VCD_WEB_PATH"},
				Destination: &cfg.Server.Path,
			},
			&cli.StringFlag{
				Name:        "output.file",
				Value:       "/etc/prometheus/vcd.json",
				Usage:       "Path to write the file_sd config",
				EnvVars:     []string{"PROMETHEUS_VCD_OUTPUT_FILE"},
				Destination: &cfg.Target.File,
			},
			&cli.IntFlag{
				Name:        "output.refresh",
				Value:       30,
				Usage:       "Discovery refresh interval in seconds",
				EnvVars:     []string{"PROMETHEUS_VCD_OUTPUT_REFRESH"},
				Destination: &cfg.Target.Refresh,
			},
			&cli.StringFlag{
				Name:    "vcd.url",
				Value:   "",
				Usage:   "URL for the vCloud Director API",
				EnvVars: []string{"PROMETHEUS_VCD_TOKEN"},
			},
			&cli.BoolFlag{
				Name:    "vcd.insecure",
				Value:   false,
				Usage:   "Accept self-signed certs for the vCloud Director API",
				EnvVars: []string{"PROMETHEUS_VCD_INSECURE"},
			},
			&cli.StringFlag{
				Name:    "vcd.username",
				Value:   "",
				Usage:   "Username for the vCloud Director API",
				EnvVars: []string{"PROMETHEUS_VCD_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "vcd.password",
				Value:   "",
				Usage:   "Password for the vCloud Director API",
				EnvVars: []string{"PROMETHEUS_VCD_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "vcd.org",
				Value:   "",
				Usage:   "Organization for the vCloud Director API",
				EnvVars: []string{"PROMETHEUS_VCD_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "vcd.config",
				Value:   "",
				Usage:   "Path to vCloud Director configuration file",
				EnvVars: []string{"PROMETHEUS_VCD_CONFIG"},
			},
		},
		Action: func(c *cli.Context) error {
			logger := setupLogger(cfg)

			if c.IsSet("vcd.config") {
				if err := readConfig(c.String("vcd.config"), cfg); err != nil {
					level.Error(logger).Log(
						"msg", "Failed to read config",
						"err", err,
					)

					return err
				}
			}

			if cfg.Target.File == "" {
				level.Error(logger).Log(
					"msg", "Missing path for output.file",
				)

				return errors.New("missing path for output.file")
			}

			if c.IsSet("vcd.url") && c.IsSet("vcd.username") && c.IsSet("vcd.password") && c.IsSet("vcd.org") {
				credentials := config.Credential{
					Project:  "default",
					URL:      c.String("vcd.url"),
					Insecure: c.Bool("vcd.insecure"),
					Username: c.String("vcd.username"),
					Password: c.String("vcd.password"),
					Org:      c.String("vcd.org"),
				}

				cfg.Target.Credentials = append(
					cfg.Target.Credentials,
					credentials,
				)

				if credentials.URL == "" {
					level.Error(logger).Log(
						"msg", "Missing required vcd.url",
					)

					return errors.New("missing required vcd.url")
				}

				if credentials.Username == "" {
					level.Error(logger).Log(
						"msg", "Missing required vcd.username",
					)

					return errors.New("missing required vcd.username")
				}

				if credentials.Password == "" {
					level.Error(logger).Log(
						"msg", "Missing required vcd.password",
					)

					return errors.New("missing required vcd.password")
				}

				if credentials.Org == "" {
					level.Error(logger).Log(
						"msg", "Missing required vcd.org",
					)

					return errors.New("missing required vcd.org")
				}
			}

			if len(cfg.Target.Credentials) == 0 {
				level.Error(logger).Log(
					"msg", "Missing any credentials",
				)

				return errors.New("missing any credentials")
			}

			return action.Server(cfg, logger)
		},
	}
}
