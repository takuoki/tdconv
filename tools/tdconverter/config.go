package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	Sheets []struct {
		Name          string `json:"name"`
		Alias         string `json:"alias"`
		SpreadsheetID string `json:"spreadsheet_id"`
	} `json:"sheets"`
}

func (c *config) AliasMap() map[string]string {
	if c == nil {
		return nil
	}
	m := map[string]string{}
	for _, s := range c.Sheets {
		m[s.Alias] = s.SpreadsheetID
	}
	return m
}

const configFile = "tdconverter.json"

type unableToReadConfigError struct {
	err error
}

func (e *unableToReadConfigError) Error() string {
	return fmt.Sprintf("Unabel to read config file (%s): ", configFile, e.err)
}

func readConfig() (*config, error) {

	s, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, &unableToReadConfigError{err: err}
	}

	conf := &config{}
	if err := json.Unmarshal(s, conf); err != nil {
		return nil, fmt.Errorf("Unabel to marshal config file (%s): %v", configFile, err)
	}

	// validate
	am := map[string]struct{}{}
	for _, s := range conf.Sheets {
		if s.SpreadsheetID == "" {
			return nil, fmt.Errorf("SpreadsheetID must not be empty (%s): %v", s.Name)
		}
		if _, ok := am[s.Alias]; ok {
			return nil, fmt.Errorf("Alias must not be duplicated (%s): %v", s.Alias)
		}
		am[s.Alias] = struct{}{}
	}

	return conf, nil
}
