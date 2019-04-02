package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
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

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Alias", "Spreadsheet ID"})
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("-")
			table.SetBorder(false)
			for _, s := range conf.Sheets {
				table.Append([]string{s.Name, s.Alias, s.SpreadsheetID})
			}
			table.Render()
			return nil
		},
	})
}
