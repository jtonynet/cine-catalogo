package responses

type RootResources struct {
	Links Links  `json:"_links"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type Links struct {
	Self Self `json:"self"`
}

type Self struct {
	HRef string `json:"href"`
}
