package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
)

type HomeConfig struct {
	ID       string `json:"id"`
	Instance string `json:"instance"`
	APIKey   string `json:"api_key"`
}

type Config struct {
	Homes     []*HomeConfig `json:"homes"`
	Remote    string        `json:"remote"`
	ExtraTags []string      `json:"extra_tags"`
}

func ParseConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c *Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}

	err = c.Valid()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Valid() error {
	var err *multierror.Error

	if len(c.Homes) == 0 {
		err = multierror.Append(err, errors.New("at least one home instance required"))
	}
	for idx, home := range c.Homes {
		errPrefix := fmt.Sprintf("home at index %d is missing its ", idx)
		if home.ID == "" {
			err = multierror.Append(err, errors.New(errPrefix+"id"))
		}
		if home.Instance == "" {
			err = multierror.Append(err, errors.New(errPrefix+"instance"))
		}
		if home.APIKey == "" {
			err = multierror.Append(err, errors.New(errPrefix+"api_key"))
		}
	}

	if c.Remote == "" {
		err = multierror.Append(err, errors.New("remote instance is missing"))
	}

	return err.ErrorOrNil()
}
