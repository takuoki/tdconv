package main

import (
	"github.com/takuoki/tdconv"
	"github.com/urfave/cli"
)

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "sql",
		Usage: "Converts the table definitions to SQL.",
		Action: func(c *cli.Context) error {
			f, err := tdconv.NewSQLFormatter()
			if err != nil {
				return err
			}
			return run(c, f)
		},
	})
}
