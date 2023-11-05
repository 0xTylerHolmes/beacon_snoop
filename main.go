package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	//flags
	remoteFlag = &cli.StringFlag{
		Name:     "remote",
		Usage:    "the beacon node to proxy e.g. http://<ip>:<rpc-port>",
		Required: true,
	}
	listenFlag = &cli.StringFlag{
		Name:        "listen",
		Usage:       "where to listen for requests",
		DefaultText: "0.0.0.0:4000",
	}
)

func main() {
	app := cli.App{
		Name:        "beacon-snoop",
		Version:     "0.1",
		Description: "snoop beacon and validator rpc",
		Flags:       []cli.Flag{remoteFlag, listenFlag},
		Action: func(c *cli.Context) error {
			config := SnooperConfig{
				remote:     c.String("remote"),
				listenAddr: c.String("listen"),
			}
			snooper, err := NewSnooper(config)
			if err != nil {
				return err
			}
			err = snooper.Snoop(context.TODO())
			if err != nil {
				return err
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
