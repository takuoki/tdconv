package main

import (
	"github.com/takuoki/tdconv"
	"github.com/urfave/cli"
)

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "go",
		Usage: "Converts the table definitions to Go struct.",
		Action: func(c *cli.Context) error {
			f, err := tdconv.NewGoFormatter()
			if err != nil {
				return err
			}
			return run(c, f)
		},
	})
}
