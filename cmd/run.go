package cmd

import (
	"github.com/cortze/ragno/crawler"

	cli "github.com/urfave/cli/v2"

	"github.com/pkg/errors"
)

var RunCommand = &cli.Command{
	Name:   "run",
	Usage:  "Run spawns an Ethereum EL crawler and starts discovering and identifying them",
	Action: RunRagno,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Define the log level of the logs it will display on the terminal",
			EnvVars:     []string{"RAGNO_LOG_LEVEL"},
			DefaultText: crawler.DefaultLogLevel,
		},
		&cli.StringFlag{
			Name:        "db-endpoint",
			Usage:       "Endpoint of the database that where the results of the crawl will be stored (needs to be initialized from before)",
			EnvVars:     []string{"RAGNO_DB_ENDPOINT"},
			DefaultText: crawler.DefaultDBEndpoint,
		},
		&cli.IntFlag{
			Name:    "disc-port",
			Usage:   "port that the tool will use for discovery purposes",
			Aliases: []string{"dp"},
			EnvVars: []string{"RAGNO_PORT"},
		},
		&cli.StringFlag{
			Name:        "ip",
			Usage:       "IP that will be assigned to the host",
			EnvVars:     []string{"RAGNO_HOST_IP"},
			DefaultText: crawler.DefaultHostIP,
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "Port that will be used by the crawler to establish TCP connections with the rest of the network",
			EnvVars: []string{"RAGNO_HOST_PORT"},
		},
		&cli.StringFlag{
			Name:        "metrics-ip",
			Usage:       "IP where the metrics of the crawler will be shown into",
			EnvVars:     []string{"RAGNO_METRICS_IP"},
			DefaultText: crawler.DefaultMetricsIP,
		},
		&cli.IntFlag{
			Name:    "metrics-port",
			Usage:   "Port that will be used to expose pprof and prometheus metrics",
			EnvVars: []string{"RAGNO_METRICS_PORT"},
		},
		&cli.StringFlag{
			Name:    "file",
			Usage:   "Path to the csv file with the Enr records to connect",
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    "concurrent-dialers",
			Usage:   "Number of workers that will be used to connect to the nodes",
			Aliases: []string{"cd"},
			EnvVars: []string{"RAGNO_DIALER_NUM"},
		},
		&cli.StringFlag{
			Name:    "concurrent-savers",
			Usage:   "Number of workers that will be used to save into the DB",
			Aliases: []string{"cs"},
			EnvVars: []string{"RAGNO_SAVER_NUM"},
		},
		&cli.IntFlag{
			Name:    "retry-amount",
			Usage:   "Number of times that the crawler will try to connect to a node before giving up",
			Aliases: []string{"ra"},
			EnvVars: []string{"RAGNO_RETRY_AMOUNT"},
		},
		&cli.IntFlag{
			Name:    "retry-delay",
			Usage:   "Number of seconds that the crawler will wait before retrying to connect to a node",
			Aliases: []string{"rd"},
			EnvVars: []string{"RAGNO_RETRY_DELAY"},
		},
	},
}

func RunRagno(ctx *cli.Context) error {
	// create a default crawler.ration
	conf := crawler.NewDefaultRun()
	err := conf.Apply(ctx)
	if err != nil {
		return errors.Wrap(err, "error applying the received configuration")
	}

	// create a new crawler from the given configuration1
	ragno, err := crawler.NewCrawler(ctx.Context, *conf)
	if err != nil {
		return errors.Wrap(err, "error initializing the crawler")
	}

	// start the crawler
	ragno.Run()

	// close the crawler
	ragno.Close()

	return nil
}
