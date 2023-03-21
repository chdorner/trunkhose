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
		h.Add("id-1", "first", uri)
		if i <= 202 {
			h.Add("id-2", "second", uri)
		}
		if i <= 262 {
			h.Add("id-1", "third", uri)
		}
	}
	require.Len(t, h.Homes["id-1"].Tags["first"].URIs, 303)
	require.Len(t, h.Homes["id-2"].Tags["second"].URIs, 202)
	require.Len(t, h.Homes["id-1"].Tags["third"].URIs, 262)

	h.Trim()

	require.Len(t, h.Homes["id-1"].Tags["first"].URIs, 250)
	require.Equal(t, "uri-303", h.Homes["id-1"].Tags["first"].URIs[249])

	require.Len(t, h.Homes["id-2"].Tags["second"].URIs, 202)
	require.Equal(t, "uri-202", h.Homes["id-2"].Tags["second"].URIs[201])

	require.Len(t, h.Homes["id-1"].Tags["third"].URIs, 250)
	require.Equal(t, "uri-262", h.Homes["id-1"].Tags["third"].URIs[249])
}
