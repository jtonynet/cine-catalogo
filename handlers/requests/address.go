package requests

type Address struct {
	Country     string `json:"country" binding:"required" hateoas:"placeholder:the country of address;required:true"`
	State       string `json:"state" binding:"required" hateoas:"placeholder:the state of address;required:true"`
	Telephone   string `json:"telephone" binding:"required" hateoas:"placeholder:the telephone of address;required:true"`
	Description string `json:"description" binding:"required" hateoas:"placeholder:the description of address;required:true"`
	PostalCode  string `json:"postalCode" binding:"required" hateoas:"placeholder:the postalCode of address;required:true"`
	Name        string `json:"name" binding:"required" hateoas:"name:name;prompt:Name;placeholder:the name of address;required:true"`
}
