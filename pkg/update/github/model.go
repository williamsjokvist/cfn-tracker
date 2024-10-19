package github

type Release struct {
	AssetsURL string `json:"assets_url"` // future
	TagName   string `json:"tag_name"`
}
