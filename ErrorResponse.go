package insightly

// ErrorResponse stores general Ridder API error response
//
type ErrorResponse struct {
	Name    string `json:"Name"`
	Message string `json:"Message"`
}
