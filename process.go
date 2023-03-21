package main

import (
	"fmt"

	"github.com/chdorner/trunkhose/history"
	"github.com/chdorner/trunkhose/mastodon"
	"github.com/hashicorp/go-multierror"
)

type ProcessContext struct {
	RemoteClient *mastodon.Client
	History      *history.History
	HomeConfig   *HomeConfig
	Config       *Config
}

func process(config *Config, histPath string) error {
	hist, err := history.NewOrParse(histPath)
	if err != nil {
		return err
	}
	defer hist.Store()

	remote, err := mastodon.NewClient(config.Remote, "")
	if err != nil {
		return err
	}

	var errs error
	processCtx := &ProcessContext{
		RemoteClient: remote,
		History:      hist,
		Config:       config,
	}
	for _, homeConfig := range config.Homes {
		processCtx.HomeConfig = homeConfig
		fmt.Printf("processing %s\n", homeConfig.ID)
		err = processHome(processCtx)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		fmt.Println()
	}

	return nil
}

func processHome(ctx *ProcessContext) error {
	home, err := mastodon.NewClient(ctx.HomeConfig.Instance, ctx.HomeConfig.APIKey)
	if err != nil {
		return err
	}

	tags, err := home.FollowedTags()
	if err != nil {
		return err
	}

	for _, name := range ctx.Config.ExtraTags {
		if name == "" {
			continue
		}
		tags = append(tags, mastodon.Tag{Name: name})
	}

	var errs *multierror.Error
	var success int
	for _, tag := range tags {
		fmt.Printf("#%s", tag.Name)
		statuses, err := ctx.RemoteClient.HashtagTimeline(tag.Name)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to get hashtag timeline for %s, err: %w", tag.Name, err))
			continue
		}

		for _, status := range statuses {
			if ctx.History.Contains(ctx.HomeConfig.ID, tag.Name, status.URI) {
				continue
			}

			err = home.Search(status.URI, true)
			if err != nil {
				errs = multierror.Append(errs, fmt.Errorf("failed to import status %s, err: %w", status.URI, err))
				fmt.Print("f")
			} else {
				ctx.History.Add(ctx.HomeConfig.ID, tag.Name, status.URI)
				success += 1
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}

	if success == 0 && errs.ErrorOrNil() != nil {
		return errs.ErrorOrNil()
	}
	return nil
}
