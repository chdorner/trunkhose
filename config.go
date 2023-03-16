package main

import (
	"errors"
	"os"

	"github.com/hashicorp/go-multierror"
)

type Config struct {
	Home      string
	APIKey    string
	Remote    string
	History   string
	ExtraTags string
}

func (c *Config) Valid() error {
	var err *multierror.Error

	if c.Home == "" {
		err = multierror.Append(err, errors.New("home instance is missing"))
	}
	if c.Remote == "" {
		err = multierror.Append(err, errors.New("remote instance is missing"))
	}
	if c.APIKey == "" {
		c.APIKey = os.Getenv("TRUNKHOSE_HOME_API_TOKEN")
	}
	if c.APIKey == "" {
		err = multierror.Append(err, errors.New("api key is missing"))
	}

	return err.ErrorOrNil()
}
