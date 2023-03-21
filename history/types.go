package history

type History struct {
	Homes map[string]*HomeHistory `json:"homes"`

	path string
}

type HomeHistory struct {
	Tags map[string]*TagHistory `json:"tags"`
}

type TagHistory struct {
	URIs []string `json:"uris"`
}
