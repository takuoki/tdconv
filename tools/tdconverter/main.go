package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/takuoki/gsheets"
	"github.com/takuoki/tdconv"
	"github.com/urfave/cli"
)

const version = "1.0.0"

var (
	cmdList = []cli.Command{}
)

func main() {

	app := cli.NewApp()
	app.Name = "tdconverter"
	app.Version = version
	app.Usage = "This tool converts the table definitions to SQL and Go struct etc."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "sheetid, i",
			Value: "",
			Usage: "spreadsheet ID of the table definitions sheet.",
		},
		cli.StringFlag{
			Name:  "sheetname, n",
			Value: "",
			Usage: "sheet name of the table definitions sheet. if not specified, all sheets in the spreadsheet.",
		},
		cli.StringFlag{
			Name:  "common, c",
			Value: "",
			Usage: "spreadsheet ID of the common columns sheet.",
		},
		cli.BoolFlag{
			Name:  "multi, m",
			Usage: "flag indicating whether to output multiple files.",
		},
	}

	app.Commands = cmdList

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func run(c *cli.Context, f tdconv.Formatter) error {

	if c.GlobalString("sheetid") == "" {
		return errors.New("Global option 'sheetid' is required")
	}

	ctx := context.Background()
	gc, err := gsheets.NewForCLI(ctx, "credentials.json")
	if err != nil {
		return fmt.Errorf("Unable to create a google sheet client. "+
			"Is 'credentials.json' present correctly?: %v", err)
	}

	var sheets []string
	if c.GlobalString("sheetname") != "" {
		sheets = []string{c.GlobalString("sheetname")}
	} else {
		sheets, err = gc.GetSheetNames(ctx, c.GlobalString("sheetid"))
		if err != nil {
			return fmt.Errorf("Unable to get all sheet names using spreadsheet id: %v", err)
		}
	}

	p, err := tdconv.NewParser()
	if err != nil {
		return fmt.Errorf("Unable to create new parser: %v", err)
	}

	if c.GlobalString("common") != "" {
		s, err := gc.GetSheet(ctx, c.GlobalString("common"), "common")
		if err != nil {
			return fmt.Errorf("Unable to get common sheet values: %v", err)
		}
		err = p.SetCommonColumns(s)
		if err != nil {
			return fmt.Errorf("Unable to parse common sheet information: %v", err)
		}
	}

	var tables []*tdconv.Table
	for _, sheetname := range sheets {
		s, err := gc.GetSheet(ctx, c.GlobalString("sheetid"), sheetname)
		if err != nil {
			return fmt.Errorf("Unable to get sheet values (sheetname=%s): %v", sheetname, err)
		}
		t, err := p.Parse(s)
		if err != nil {
			return fmt.Errorf("Unable to parse sheet information (sheetname=%s): %v", sheetname, err)
		}
		tables = append(tables, t)
	}

	outdir := "./out/" + c.Command.Name
	if _, err := os.Stat("./out"); err != nil {
		os.Mkdir("./out", 0777)
	}
	if _, err := os.Stat(outdir); err != nil {
		os.Mkdir(outdir, 0777)
	}

	if err = tdconv.Output(f, tables, c.GlobalBool("multi"), outdir, c.Command.Name); err != nil {
		return fmt.Errorf("Fail to output table definitions: %v", err)
	}

	fmt.Println("complete!")

	return nil
}
