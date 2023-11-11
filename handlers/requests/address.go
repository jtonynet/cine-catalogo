package requests

type Address struct {
	Country     string `json:"country" binding:"required"`
	State       string `json:"state" binding:"required"`
	Telephone   string `json:"telephone" binding:"required"`
	Description string `json:"description" binding:"required"`
	PostalCode  string `json:"postalCode" binding:"required"`
	Name        string `json:"name" binding:"required"`
}
