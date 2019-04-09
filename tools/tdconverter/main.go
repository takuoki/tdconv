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
			Value: "1MWfimYqzTtHwuw4i8JCZZwDnsvLCBVQGiOyMpH8-2IQ",
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

	err := validate(c)
	if err != nil {
		return err
	}

	am, err := getAliasMap()
	if err != nil {
		return err
	}

	ctx := context.Background()

	gc, err := gsheets.NewForCLI(ctx, "credentials.json")
	if err != nil {
		return fmt.Errorf("Unable to create a google sheet client. "+
			"Is 'credentials.json' present correctly?: %v", err)
	}

	sheetid := c.GlobalString("sheetid")
	if s, ok := am[sheetid]; ok {
		sheetid = s
	}

	common := c.GlobalString("common")
	if s, ok := am[common]; ok {
		common = s
	}

	ts, err := parse(ctx, gc, sheetid, c.GlobalString("sheetname"), common)
	if err != nil {
		return err
	}

	err = output(f, c.Command.Name, ts, c.GlobalBool("multi"))
	if err != nil {
		return err
	}

	fmt.Println("complete!")

	return nil
}

func validate(c *cli.Context) error {

	if c.GlobalString("sheetid") == "" {
		return errors.New("Global option 'sheetid' is required")
	}

	return nil
}

func getAliasMap() (map[string]string, error) {

	conf, err := readConfig()
	if err != nil {
		switch err.(type) {
		case *unableToReadConfigError:
			return nil, nil
		default:
			return nil, err
		}
	}

	return conf.AliasMap(), nil
}

func parse(ctx context.Context, gc *gsheets.Client, id, sheet, common string) (*tdconv.TableSet, error) {

	title, err := gc.GetTitle(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Unable to get spreadsheet title: %v", err)
	}

	var sheets []string
	if sheet != "" {
		sheets = []string{sheet}
	} else {
		sheets, err = gc.GetSheetNames(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("Unable to get all sheet names using spreadsheet id: %v", err)
		}
	}

	p, err := tdconv.NewParser()
	if err != nil {
		return nil, fmt.Errorf("Unable to create new parser: %v", err)
	}

	if common != "" {
		s, err := gc.GetSheet(ctx, common, "common")
		if err != nil {
			return nil, fmt.Errorf("Unable to get common sheet values: %v", err)
		}
		err = p.SetCommonColumns(s)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse common sheet information: %v", err)
		}
	}

	var tables []*tdconv.Table
	for _, sheetname := range sheets {
		s, err := gc.GetSheet(ctx, id, sheetname)
		if err != nil {
			return nil, fmt.Errorf("Unable to get sheet values (sheetname=%s): %v", sheetname, err)
		}
		t, err := p.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse sheet information (sheetname=%s): %v", sheetname, err)
		}
		tables = append(tables, t)
	}

	return &tdconv.TableSet{
		Name:   title,
		Tables: tables,
	}, nil
}

func output(f tdconv.Formatter, commandName string, ts *tdconv.TableSet, multi bool) error {

	outdir := "./out/" + commandName
	if _, err := os.Stat("./out"); os.IsNotExist(err) {
		os.Mkdir("./out", 0777)
	}
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.Mkdir(outdir, 0777)
	}

	if err := tdconv.Output(f, ts, multi, outdir); err != nil {
		return fmt.Errorf("Fail to output table definitions: %v", err)
	}

	return nil
}
