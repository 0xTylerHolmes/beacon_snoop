package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	remote      string
	listen      string
	logHeaders  bool
	logToFile   bool
	logFilePath string
	//flags
	remoteFlag = &cli.StringFlag{
		Name:        "remote",
		Usage:       "the beacon node to proxy e.g. http://<ip>:<rpc-port>",
		Required:    true,
		Destination: &remote,
	}
	listenFlag = &cli.StringFlag{
		Name:        "listen",
		Usage:       "where to listen for requests",
		DefaultText: ":4000",
		Required:    true,
		Destination: &listen,
	}
	logHeadersFlag = &cli.BoolFlag{
		Name:        "logHeaders",
		Usage:       "if the file logger should log the headers for requests/responses",
		DefaultText: "false",
		Destination: &logHeaders,
	}
	logToFileFlag = &cli.BoolFlag{
		Name:        "logToFile",
		Usage:       "whether or not to log verbose information to file",
		DefaultText: "false",
		Destination: &logToFile,
	}

	logFileFlag = &cli.PathFlag{
		Name:        "logFile",
		Usage:       "where to log verbose information to file. If not supplied then don't log to file",
		DefaultText: "/tmp/beacon_snoop.log",
		Destination: &logFilePath,
	}
)

func main() {
	app := cli.App{
		Name:        "beacon-snoop",
		Version:     "0.1",
		Description: "snoop beacon and validator rpc",
		Flags:       []cli.Flag{remoteFlag, listenFlag, logHeadersFlag, logFileFlag, logToFileFlag},
		Action: func(c *cli.Context) error {
			config := SnooperConfig{
				remote:      remote,
				listenAddr:  listen,
				logHeaders:  logHeaders,
				logFilePath: logFilePath,
				logToFile:   logToFile,
			}
			snooper, err := NewSnooper(config)
			if err != nil {
				return err
			}
			err = snooper.Snoop(context.Background())
			if err != nil {
				return err
			}
			return nil
		},
		UsageText: "./snooper --listen :4000 --remote http://10.0.20.4:5052 --logHeaders --logFile /tmp/beacon_snooper_log.log",
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
