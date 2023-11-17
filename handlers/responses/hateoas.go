package responses

type HATEOASResult struct {
	Embedded  interface{} `json:"_embedded"`
	Links     interface{} `json:"_links"`
	Templates interface{} `json:"_templates"`
}
type HATEOASListItemProperties struct {
	Links interface{} `json:"_links"`
}

type HREFObject struct {
	HREF string `json:"href"`
}

type HATEOASTemplateParams struct {
	Name          string
	ResourceURL   string
	HTTPMethod    string
	RequestStruct interface{}
}
