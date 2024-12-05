package main

import (
	"context"
	"log"
	"os"
	"sort"

	"github.com/lordvidex/gostream/cmd/gostream/client"
	"github.com/lordvidex/gostream/cmd/gostream/server"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:    "gostream",
		Usage:   "gostream CLI tool",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Aliases:   []string{"c"},
				TakesFile: true,
			},
		},
		Commands: []*cli.Command{server.Cmd, client.Cmd},
	}

	sort.Sort(cli.FlagsByName(cmd.Flags))

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
