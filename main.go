package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/chdorner/trunkhose/history"
	"github.com/chdorner/trunkhose/mastodon"
	"github.com/hashicorp/go-multierror"
	flag "github.com/spf13/pflag"
)

func main() {
	var config Config
	flag.StringVar(&config.Home, "home", "", "hostname of the home instance")
	flag.StringVar(&config.APIKey, "api-key", "", "api key for the home instance")
	flag.StringVar(&config.Remote, "remote", "", "host of the remote instance")
	flag.StringVar(&config.History, "history", "", "path to the history file")
	flag.StringVar(&config.ExtraTags, "extra-tags", "", "extra tags to search")
	flag.Parse()

	err := config.Valid()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "see --help for more information.")
		os.Exit(1)
	}

	if config.History == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		config.History = path.Join(wd, "trunkhose-history.json")
	}

	err = herd(&config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func herd(config *Config) error {
	hist, err := history.NewOrParse(config.History)
	if err != nil {
		return err
	}
	defer hist.Store()

	home, err := mastodon.NewClient(config.Home, config.APIKey)
	if err != nil {
		return err
	}
	remote, err := mastodon.NewClient(config.Remote, "")
	if err != nil {
		return err
	}

	tags, err := home.FollowedTags()
	if err != nil {
		return err
	}

	for _, name := range strings.Split(config.ExtraTags, ",") {
		tags = append(tags, mastodon.Tag{Name: name})
	}

	var errs *multierror.Error
	for _, tag := range tags {
		fmt.Printf("#%s", tag.Name)
		statuses, err := remote.HashtagTimeline(tag.Name)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		for _, status := range statuses {
			if hist.Contains(tag.Name, status.URI) {
				continue
			}

			err = home.Search(status.URI, true)
			if err != nil {
				errs = multierror.Append(errs, err)
				fmt.Print("f")
			} else {
				hist.Add(tag.Name, status.URI)
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}

	return errs.ErrorOrNil()
}
