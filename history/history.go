package history

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

const (
	Limit = 250
)

func NewOrParse(path string) (*History, error) {
	_, err := os.Stat(path)
	if err == nil {
		return Parse(path)
	}
	return New(path), nil
}

func New(path string) *History {
	return &History{
		Homes: make(map[string]*HomeHistory),
		path:  path,
	}
}

func Parse(path string) (*History, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var h *History
	err = json.NewDecoder(f).Decode(&h)
	if err != nil {
		return nil, err
	}
	h.path = path

	return h, nil
}

func (h *History) Contains(id, tag, uri string) bool {
	home, ok := h.Homes[id]
	if !ok {
		return false
	}

	tagHist, ok := home.Tags[tag]
	if !ok {
		return false
	}

	return slices.Contains(tagHist.URIs, uri)
}

func (h *History) Add(id, tag, uri string) {
	if h.Contains(id, tag, uri) {
		return
	}

	if _, ok := h.Homes[id]; !ok {
		h.Homes[id] = &HomeHistory{
			Tags: make(map[string]*TagHistory),
		}
	}
	home := h.Homes[id]

	if _, ok := home.Tags[tag]; !ok {
		home.Tags[tag] = &TagHistory{}
	}
	tagHist := home.Tags[tag]
	tagHist.URIs = append(tagHist.URIs, uri)
}

func (h *History) Store() error {
	f, err := os.OpenFile(h.path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write history: %s", err)
		return err
	}
	defer f.Close()

	h.Trim()
	err = json.NewEncoder(f).Encode(h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write history: %s", err)
		return err
	}
	return nil
}

func (h *History) Trim() {
	for _, homeHist := range h.Homes {
		for _, tagHist := range homeHist.Tags {
			length := len(tagHist.URIs)
			if length <= Limit {
				continue
			}
			tagHist.URIs = tagHist.URIs[length-Limit : len(tagHist.URIs)]
		}
	}
}
