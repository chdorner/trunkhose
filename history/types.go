package history

type History struct {
	Tags map[string]*TagHistory `json:"tags"`

	path string
}

type TagHistory struct {
	URIs []string `json:"uris"`
}
