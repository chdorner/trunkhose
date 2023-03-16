package mastodon

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	Host   string
	APIKey string

	client *http.Client
}

func NewClient(host, apiKey string) (*Client, error) {
	if host == "" {
		return nil, errors.New("host is required")
	}
	return &Client{
		host,
		apiKey,
		&http.Client{},
	}, nil
}

func (mc *Client) FollowedTags() ([]Tag, error) {
	resp, err := mc.get("/api/v1/followed_tags", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tags []Tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (mc *Client) Search(query string, resolve bool) error {
	q := url.Values{}
	q.Add("q", query)
	q.Add("resolve", fmt.Sprintf("%v", resolve))
	resp, err := mc.get("/api/v2/search", q)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (mc *Client) HashtagTimeline(tag string) ([]Status, error) {
	q := url.Values{}
	q.Add("limit", "40")
	url := fmt.Sprintf("/api/v1/timelines/tag/%s", tag)
	resp, err := mc.get(url, q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var statuses []Status
	err = json.NewDecoder(resp.Body).Decode(&statuses)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (mc *Client) get(endpoint string, query url.Values) (*http.Response, error) {
	url := fmt.Sprintf("https://%s%s", mc.Host, endpoint)
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if mc.APIKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mc.APIKey))
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	return mc.client.Do(req)
}
