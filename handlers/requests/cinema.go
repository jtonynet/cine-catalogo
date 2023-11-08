package requests

type Cinema struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Capacity    int64  `json:"capacity" binding:"required"`
}
