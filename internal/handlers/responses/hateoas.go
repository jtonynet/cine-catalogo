package responses

type HATEOASListResult struct {
	Embedded  interface{} `json:"_embedded"`
	Links     interface{} `json:"_links"`
	Templates interface{} `json:"_templates"`
}
type HATEOASListItemResult struct {
	Links     interface{} `json:"_links"`
	Templates interface{} `json:"_templates,omitempty"`
}

type HATEOASLink struct {
	HREF string `json:"href"`
}
