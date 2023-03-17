package history_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/chdorner/trunkhose/history"
	"github.com/stretchr/testify/require"
)

func TestTrim(t *testing.T) {
	h := history.New(path.Join("/", "tmp", "trunkhose-history.json"))
	for i := 1; i <= 303; i++ {
		uri := fmt.Sprintf("uri-%d", i)
		h.Add("first", uri)
		if i <= 202 {
			h.Add("second", uri)
		}
		if i <= 262 {
			h.Add("third", uri)
		}
	}
	require.Len(t, h.Tags["first"].URIs, 303)
	require.Len(t, h.Tags["second"].URIs, 202)
	require.Len(t, h.Tags["third"].URIs, 262)

	h.Trim()

	require.Len(t, h.Tags["first"].URIs, 250)
	require.Equal(t, "uri-303", h.Tags["first"].URIs[249])

	require.Len(t, h.Tags["second"].URIs, 202)
	require.Equal(t, "uri-202", h.Tags["second"].URIs[201])

	require.Len(t, h.Tags["third"].URIs, 250)
	require.Equal(t, "uri-262", h.Tags["third"].URIs[249])
}
