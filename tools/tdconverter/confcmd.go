package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "conf",
		Usage: fmt.Sprintf("Validate and show configuration file (%s).", configFile),
		Action: func(c *cli.Context) error {

			conf, err := readConfig()
			if err != nil {
				switch err.(type) {
				case *unableToReadConfigError:
					fmt.Printf("There are no configuration file (%s).\n", configFile)
					return nil
				default:
					return err
				}
			}

			fmt.Println("Alias:\tSpreadsheetID")
			fmt.Println("--------------------")
			for _, s := range conf.Sheets {
				fmt.Printf("%s:\t%s\n", s.Alias, s.SpreadsheetID)
			}
			return nil
		},
	})
}
